package crd

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/to"
	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apiyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

type Config struct {
	Logger micrologger.Logger

	ApiextensionsReference string
	GitHubToken            string
	Provider               string
}

type CRDGetter struct {
	logger micrologger.Logger

	apiextensionsReference string
	githubClient           *github.Client
	provider               string
}

var (
	crdV1GVK = schema.GroupVersionKind{
		Group:   "apiextensions.k8s.io",
		Version: "v1",
		Kind:    "CustomResourceDefinition",
	}
	crdV1Beta1GVK = schema.GroupVersionKind{
		Group:   "apiextensions.k8s.io",
		Version: "v1beta1",
		Kind:    "CustomResourceDefinition",
	}
)

func NewCRDGetter(config Config) (*CRDGetter, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	var githubClient *github.Client
	{
		var tc *http.Client
		if config.GitHubToken != "" {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: config.GitHubToken},
			)
			tc = oauth2.NewClient(context.Background(), ts)
		} else {
			tc = http.DefaultClient
		}

		githubClient = github.NewClient(tc)
	}

	crdGetter := &CRDGetter{
		logger: config.Logger,

		apiextensionsReference: config.ApiextensionsReference,
		githubClient:           githubClient,
		provider:               config.Provider,
	}

	return crdGetter, nil
}

func (g CRDGetter) LoadCRD(ctx context.Context, group, kind string) (*apiextensionsv1.CustomResourceDefinition, error) {
	var crds []*apiextensionsv1.CustomResourceDefinition
	charts := []string{
		"crds-common",
	}

	if g.provider != "" {
		charts = append(charts, fmt.Sprintf("cds-%s", g.provider))
	}

	for _, chart := range charts {
		chartCRDs, err := downloadHelmChartCRDs(ctx, g.githubClient, chart, g.apiextensionsReference)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		crds = append(crds, chartCRDs...)
	}

	for _, crd := range crds {
		if crd.Spec.Names.Kind == kind && crd.Spec.Group == group {
			return crd, nil
		}
	}

	return nil, microerror.Maskf(notFoundError, "CRD kind %#q not found in group %#q", kind, group)
}

func convertCRDV1Beta1(original *apiextensionsv1beta1.CustomResourceDefinition) (*apiextensionsv1.CustomResourceDefinition, error) {
	var hub apiextensions.CustomResourceDefinition
	err := apiextensionsv1beta1.Convert_v1beta1_CustomResourceDefinition_To_apiextensions_CustomResourceDefinition(original, &hub, nil)
	if err != nil {
		return nil, err
	}

	var spoke apiextensionsv1.CustomResourceDefinition
	err = apiextensionsv1.Convert_apiextensions_CustomResourceDefinition_To_v1_CustomResourceDefinition(&hub, &spoke, nil)
	if err != nil {
		return nil, err
	}

	if spoke.Spec.Versions[0].Schema == nil {
		spoke.Spec.Versions[0].Schema = &apiextensionsv1.CustomResourceValidation{
			OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
				Type:                   "object",
				XPreserveUnknownFields: to.BoolP(true),
			},
		}
	}

	return &spoke, nil
}

func decodeCRDs(readCloser io.ReadCloser) ([]*apiextensionsv1.CustomResourceDefinition, error) {
	reader := apiyaml.NewYAMLReader(bufio.NewReader(readCloser))
	decoder := scheme.Codecs.UniversalDecoder()

	defer func(contentReader io.ReadCloser) {
		err := readCloser.Close()
		if err != nil {
			panic(err)
		}
	}(readCloser)

	var crds []*apiextensionsv1.CustomResourceDefinition

	for {
		doc, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, microerror.Mask(err)
		}

		//  Skip over empty documents, i.e. a leading `---`
		if len(bytes.TrimSpace(doc)) == 0 {
			continue
		}

		var object unstructured.Unstructured
		_, decodedGVK, err := decoder.Decode(doc, nil, &object)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		switch *decodedGVK {
		case crdV1GVK:
			var crd apiextensionsv1.CustomResourceDefinition
			_, _, err = decoder.Decode(doc, nil, &crd)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			crds = append(crds, &crd)
		case crdV1Beta1GVK:
			var crd apiextensionsv1beta1.CustomResourceDefinition
			_, _, err = decoder.Decode(doc, nil, &crd)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			converted, err := convertCRDV1Beta1(&crd)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			crds = append(crds, converted)
		}
	}

	return crds, nil
}

func downloadHelmChartCRDs(ctx context.Context, client *github.Client, helmChart string, ref string) ([]*apiextensionsv1.CustomResourceDefinition, error) {
	getOptions := github.RepositoryContentGetOptions{
		Ref: ref,
	}

	templatesPath := path.Join("helm", helmChart, "templates")
	_, contents, _, err := client.Repositories.GetContents(ctx, "giantswarm", "apiextensions", templatesPath, &getOptions)
	if err != nil {
		return nil, err
	}

	var allCrds []*apiextensionsv1.CustomResourceDefinition
	for _, file := range contents {
		filePath := path.Join(templatesPath, *file.Name)
		contentReader, _, err := client.Repositories.DownloadContents(ctx, "giantswarm", "apiextensions", filePath, &getOptions)
		if err != nil {
			return nil, err
		}

		crds, err := decodeCRDs(contentReader)
		if err != nil {
			return nil, err
		}

		allCrds = append(allCrds, crds...)
	}

	return allCrds, nil
}
