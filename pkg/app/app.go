package app

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/giantswarm/apiextensions/v3/pkg/annotation"
	applicationv1alpha1 "github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/apiextensions/v3/pkg/label"
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type Config struct {
	AppCatalog          string
	AppName             string
	AppNamespace        string
	AppVersion          string
	ConfigVersion       string
	DisableForceUpgrade bool
	Name                string
	UserConfigMapName   string
	UserSecretName      string
}

// NewCR returns new application CR.
//
// AppCatalog is the name of the app catalog where the app stored.
func NewCR(c Config) *applicationv1alpha1.App {
	annotations := map[string]string{}
	{
		if c.ConfigVersion != "" {
			annotations[annotation.ConfigVersion] = c.ConfigVersion
		}
		if !c.DisableForceUpgrade {
			annotations["chart-operator.giantswarm.io/force-helm-upgrade"] = "true"
		}
	}

	var userConfig applicationv1alpha1.AppSpecUserConfig
	if c.UserConfigMapName != "" {
		userConfig.ConfigMap = applicationv1alpha1.AppSpecUserConfigConfigMap{
			Name:      c.UserConfigMapName,
			Namespace: "giantswarm",
		}
	}
	if c.UserSecretName != "" {
		userConfig.Secret = applicationv1alpha1.AppSpecUserConfigSecret{
			Name:      c.UserSecretName,
			Namespace: "giantswarm",
		}
	}

	appCR := &applicationv1alpha1.App{
		TypeMeta: applicationv1alpha1.NewAppTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:        c.Name,
			Namespace:   "giantswarm",
			Annotations: annotations,
			Labels: map[string]string{
				// Version 0.0.0 means this is reconciled by
				// unique operator.
				label.AppOperatorVersion: "0.0.0",
			},
		},
		Spec: applicationv1alpha1.AppSpec{
			Catalog: c.AppCatalog,
			KubeConfig: applicationv1alpha1.AppSpecKubeConfig{
				InCluster: true,
			},
			Name:       c.AppName,
			Namespace:  c.AppNamespace,
			Version:    c.AppVersion,
			UserConfig: userConfig,
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
