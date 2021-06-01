package values

import (
	"context"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/imdario/mergo"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

// Config represents the configuration used to create a new values service.
type Config struct {
	// Dependencies.
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
}

// Values implements the values service.
type Values struct {
	// Dependencies.
	k8sClient kubernetes.Interface
	logger    micrologger.Logger
}

// New creates a new configured values service.
func New(config Config) (*Values, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &Values{
		// Dependencies.
		k8sClient: config.K8sClient,
		logger:    config.Logger,
	}

	return r, nil
}

// MergeAll merges both configmap and secret values to produce a single set of
// values that can be passed to Helm.
func (v *Values) MergeAll(ctx context.Context, app v1alpha1.App, catalog v1alpha1.Catalog) (map[string]interface{}, error) {
	configMapData, err := v.MergeConfigMapData(ctx, app, catalog)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secretData, err := v.MergeSecretData(ctx, app, catalog)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	err = mergo.Merge(&configMapData, secretData, mergo.WithOverride)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return configMapData, nil
}

func extractData(resourceType, name string, data map[string]string) (map[string]interface{}, error) {
	var err error
	var rawMapData map[string]interface{}

	if data == nil {
		return rawMapData, nil
	}

	if len(data) != 1 {
		return nil, microerror.Maskf(parsingError, "expected %#q %s has only one key but got %d", name, resourceType, len(data))
	}

	var rawData []byte
	for _, v := range data {
		rawData = []byte(v)
	}

	err = yaml.Unmarshal(rawData, &rawMapData)
	if err != nil {
		return nil, microerror.Maskf(parsingError, "failed to parse %#q %s, logs: %s", name, resourceType, err.Error())
	}

	return rawMapData, nil
}

// toStringMap converts from a byte slice map to a string map.
func toStringMap(input map[string][]byte) map[string]string {
	if input == nil {
		return nil
	}

	result := map[string]string{}

	for k, v := range input {
		result[k] = string(v)
	}

	return result
}
