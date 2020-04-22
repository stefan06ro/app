package app

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ghodss/yaml"
	applicationv1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	AppCatalog          string
	AppName             string
	AppNamespace        string
	AppVersion          string
	DisableForceUpgrade bool
	Name                string
}

// NewCR returns new application CR.
//
// AppCatalog is the name of the app catalog where the app stored.
func NewCR(c Config) *applicationv1alpha1.App {
	var annotations map[string]string
	{
		if !c.DisableForceUpgrade {
			annotations = map[string]string{
				"chart-operator.giantswarm.io/force-helm-upgrade": "true",
			}
		}
	}

	appCR := &applicationv1alpha1.App{
		TypeMeta: applicationv1alpha1.NewAppTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:        c.Name,
			Namespace:   "giantswarm",
			Annotations: annotations,
			Labels: map[string]string{
				"app-operator.giantswarm.io/version": "1.0.0",
			},
		},
		Spec: applicationv1alpha1.AppSpec{
			Catalog: c.AppCatalog,
			KubeConfig: applicationv1alpha1.AppSpecKubeConfig{
				InCluster: true,
			},
			Name:      c.AppName,
			Namespace: c.AppNamespace,
			Version:   c.AppVersion,
		},
	}

	return appCR
}

func Marshal(appCR *applicationv1alpha1.App, format string) (string, error) {
	var output []byte
	var err error

	switch format {
	case "json":
		output, err = json.Marshal(appCR)
		if err != nil {
			return "", microerror.Mask(err)
		}
	case "yaml":
		output, err = yaml.Marshal(appCR)
		if err != nil {
			return "", microerror.Mask(err)
		}
	default:
		return "", microerror.Maskf(executionFailedError, "format: %q", format)
	}

	return string(output), nil
}

func Print(w io.Writer, format string, appCR *applicationv1alpha1.App) error {
	output, err := Marshal(appCR, format)
	if err != nil {
		return microerror.Mask(err)
	}

	_, err = fmt.Fprintf(w, "%s", output)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
