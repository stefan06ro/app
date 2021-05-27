package key

import (
	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/microerror"
)

func CatalogTitle(customResource v1alpha1.Catalog) string {
	return customResource.Spec.Title
}

func CatalogStorageURL(customResource v1alpha1.Catalog) string {
	return customResource.Spec.Storage.URL
}

func CatalogConfigMapName(customResource v1alpha1.Catalog) string {
	config := catalogConfig(customResource)

	if config != nil && config.ConfigMap != nil {
		return customResource.Spec.Config.ConfigMap.Name
	}

	return ""
}

func CatalogConfigMapNamespace(customResource v1alpha1.Catalog) string {
	config := catalogConfig(customResource)

	if config != nil && config.ConfigMap != nil {
		return customResource.Spec.Config.ConfigMap.Namespace
	}

	return ""
}

func CatalogSecretName(customResource v1alpha1.Catalog) string {
	config := catalogConfig(customResource)

	if config != nil && config.Secret != nil {
		return customResource.Spec.Config.Secret.Name
	}

	return ""
}

func CatalogSecretNamespace(customResource v1alpha1.Catalog) string {
	config := catalogConfig(customResource)

	if config != nil && config.Secret != nil {
		return customResource.Spec.Config.Secret.Namespace
	}

	return ""
}

func CatalogType(customResource v1alpha1.Catalog) string {
	if val, ok := customResource.ObjectMeta.Labels[label.CatalogType]; ok {
		return val
	}

	return ""
}

func CatalogVisibility(customResource v1alpha1.Catalog) string {
	if val, ok := customResource.ObjectMeta.Labels[label.CatalogVisibility]; ok {
		return val
	}

	return ""
}

func ToCatalog(v interface{}) (v1alpha1.Catalog, error) {
	customResource, ok := v.(*v1alpha1.Catalog)
	if !ok {
		return v1alpha1.Catalog{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha1.Catalog{}, v)
	}

	if customResource == nil {
		return v1alpha1.Catalog{}, microerror.Maskf(emptyValueError, "empty value cannot be converted to CustomObject")
	}

	return *customResource, nil
}

func catalogConfig(customResource v1alpha1.Catalog) *v1alpha1.CatalogSpecConfig {
	return customResource.Spec.Config
}
