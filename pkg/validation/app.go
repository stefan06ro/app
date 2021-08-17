package validation

import (
	"context"
	"fmt"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/app/v5/pkg/key"
)

const (
	catalogNotFoundTemplate         = "catalog %#q not found"
	nameTooLongTemplate             = "name %#q is %d chars and exceeds max length of %d chars"
	namespaceNotFoundReasonTemplate = "namespace is not specified for %s %#q"
	labelInvalidValueTemplate       = "label %#q has invalid value %#q"
	labelNotFoundTemplate           = "label %#q not found"
	resourceNotFoundTemplate        = "%s %#q in namespace %#q not found"

	defaultCatalogName = "default"

	// nameMaxLength is 53 characters as this is the maximum allowed for Helm
	// release names.
	nameMaxLength = 53
)

func (v *Validator) ValidateApp(ctx context.Context, app v1alpha1.App) (bool, error) {
	var err error

	err = v.validateCatalog(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateConfig(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateKubeConfig(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateLabels(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateMetadataConstraints(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateName(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateNamespaceConfig(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateUserConfig(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	return true, nil
}

func (v *Validator) validateCatalog(ctx context.Context, cr v1alpha1.App) error {
	var err error

	if key.CatalogName(cr) == "" {
		return nil
	}

	var namespaces []string
	{
		if key.CatalogNamespace(cr) != "" {
			namespaces = []string{key.CatalogNamespace(cr)}
		} else {
			namespaces = []string{metav1.NamespaceDefault, "giantswarm"}
		}
	}

	var catalog *v1alpha1.Catalog

	for _, ns := range namespaces {
		catalog, err = v.g8sClient.ApplicationV1alpha1().Catalogs(ns).Get(ctx, key.CatalogName(cr), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			// no-op
			continue
		} else if err != nil {
			return microerror.Mask(err)
		}
		break
	}

	if catalog == nil || catalog.Name == "" {
		return microerror.Maskf(validationError, catalogNotFoundTemplate, key.CatalogName(cr))
	}

	return nil
}

func (v *Validator) validateConfig(ctx context.Context, cr v1alpha1.App) error {
	if key.AppConfigMapName(cr) != "" {
		ns := key.AppConfigMapNamespace(cr)
		if ns == "" {
			return microerror.Maskf(validationError, namespaceNotFoundReasonTemplate, "configmap", key.AppConfigMapName(cr))
		}

		_, err := v.k8sClient.CoreV1().ConfigMaps(ns).Get(ctx, key.AppConfigMapName(cr), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			// appConfigMapNotFoundError is used rather than a validation error because
			// during cluster creation there is a short delay while it is generated.
			return microerror.Maskf(appConfigMapNotFoundError, resourceNotFoundTemplate, "configmap", key.AppConfigMapName(cr), ns)
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	if key.AppSecretName(cr) != "" {
		ns := key.AppSecretNamespace(cr)
		if ns == "" {
			return microerror.Maskf(validationError, namespaceNotFoundReasonTemplate, "secret", key.AppSecretName(cr))
		}

		_, err := v.k8sClient.CoreV1().Secrets(ns).Get(ctx, key.AppSecretName(cr), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			return microerror.Maskf(validationError, resourceNotFoundTemplate, "secret", key.AppSecretName(cr), ns)
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (v *Validator) validateName(ctx context.Context, cr v1alpha1.App) error {
	if len(cr.Name) > nameMaxLength {
		return microerror.Maskf(validationError, nameTooLongTemplate, cr.Name, len(cr.Name), nameMaxLength)
	}

	return nil
}

func (v *Validator) validateNamespaceConfig(ctx context.Context, cr v1alpha1.App) error {
	annotations := key.AppNamespaceAnnotations(cr)
	labels := key.AppNamespaceLabels(cr)

	if annotations == nil && labels == nil {
		// no-op
		return nil
	}

	var apps []v1alpha1.App
	{
		lo := metav1.ListOptions{
			FieldSelector: fmt.Sprintf("metadata.name!=%s", cr.Name),
		}
		appList, err := v.g8sClient.ApplicationV1alpha1().Apps(cr.Namespace).List(ctx, lo)
		if err != nil {
			return microerror.Mask(err)
		}

		apps = appList.Items
	}

	for _, app := range apps {
		if key.AppNamespace(cr) != key.AppNamespace(app) {
			continue
		}

		targetAnnotations := key.AppNamespaceAnnotations(app)
		if targetAnnotations != nil && annotations != nil {
			for k, v := range targetAnnotations {
				originalValue, ok := annotations[k]
				if ok && originalValue != v {
					return microerror.Maskf(validationError, "app %#q annotation %#q for target namespace %#q collides with value %#q for app %#q",
						key.AppName(cr), k, key.AppNamespace(cr), v, app.Name)
				}
			}
		}

		targetLabels := key.AppNamespaceLabels(app)
		if targetLabels != nil && labels != nil {
			for k, v := range targetLabels {
				originalValue, ok := labels[k]
				if ok && originalValue != v {
					return microerror.Maskf(validationError, "app %#q label %#q for target namespace %#q collides with value %#q for app %#q",
						key.AppName(cr), k, key.AppNamespace(cr), v, app.Name)
				}
			}
		}
	}

	return nil
}

func (v *Validator) validateKubeConfig(ctx context.Context, cr v1alpha1.App) error {
	if !key.InCluster(cr) {
		ns := key.KubeConfigSecretNamespace(cr)
		if ns == "" {
			return microerror.Maskf(validationError, namespaceNotFoundReasonTemplate, "kubeconfig secret", key.KubeConfigSecretName(cr))
		}

		_, err := v.k8sClient.CoreV1().Secrets(key.KubeConfigSecretNamespace(cr)).Get(ctx, key.KubeConfigSecretName(cr), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			// kubeConfigNotFoundError is used rather than a validation error because
			// during cluster creation there is a short delay while it is generated.
			return microerror.Maskf(kubeConfigNotFoundError, resourceNotFoundTemplate, "kubeconfig secret", key.KubeConfigSecretName(cr), ns)
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (v *Validator) validateLabels(ctx context.Context, cr v1alpha1.App) error {
	if key.VersionLabel(cr) == "" {
		return microerror.Maskf(validationError, labelNotFoundTemplate, label.AppOperatorVersion)
	}
	if key.VersionLabel(cr) == key.LegacyAppVersionLabel {
		return microerror.Maskf(validationError, labelInvalidValueTemplate, label.AppOperatorVersion, key.VersionLabel(cr))
	}

	return nil
}

func (v *Validator) validateMetadataConstraints(ctx context.Context, cr v1alpha1.App) error {
	name := key.AppCatalogEntryName(key.CatalogName(cr), key.AppName(cr), key.Version(cr))

	entry, err := v.g8sClient.ApplicationV1alpha1().AppCatalogEntries(metav1.NamespaceDefault).Get(ctx, name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		v.logger.Debugf(ctx, "appcatalogentry %#q not found, skipping metadata validation", name)
		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	if entry.Spec.Restrictions == nil {
		// no-op
		return nil
	}

	if len(entry.Spec.Restrictions.CompatibleProviders) > 0 {
		if !contains(entry.Spec.Restrictions.CompatibleProviders, v1alpha1.Provider(v.provider)) {
			return microerror.Maskf(validationError, "app %#q can only be installed for providers %#q not %#q",
				cr.Spec.Name, entry.Spec.Restrictions.CompatibleProviders, v.provider)
		}
	}

	if entry.Spec.Restrictions.FixedNamespace != "" {
		if entry.Spec.Restrictions.FixedNamespace != cr.Spec.Namespace {
			return microerror.Maskf(validationError, "app %#q can only be installed in namespace %#q only, not %#q",
				cr.Spec.Name, entry.Spec.Restrictions.FixedNamespace, cr.Spec.Namespace)
		}
	}

	var apps []v1alpha1.App
	if entry.Spec.Restrictions.ClusterSingleton || entry.Spec.Restrictions.NamespaceSingleton {
		lo := metav1.ListOptions{
			FieldSelector: fmt.Sprintf("metadata.name!=%s", cr.Name),
		}
		appList, err := v.g8sClient.ApplicationV1alpha1().Apps(cr.Namespace).List(ctx, lo)
		if err != nil {
			return microerror.Mask(err)
		}

		apps = appList.Items
	}

	for _, app := range apps {
		if app.Spec.Name == cr.Spec.Name {
			if entry.Spec.Restrictions.ClusterSingleton {
				return microerror.Maskf(validationError, "app %#q can only be installed once in cluster %#q",
					cr.Spec.Name, cr.Namespace)
			}
			if entry.Spec.Restrictions.NamespaceSingleton {
				if app.Spec.Namespace == cr.Spec.Namespace {
					return microerror.Maskf(validationError, "app %#q can only be installed only once in namespace %#q",
						cr.Spec.Name, key.Namespace(cr))
				}
			}
		}
	}

	return nil
}

func (v *Validator) validateUserConfig(ctx context.Context, cr v1alpha1.App) error {
	if key.UserConfigMapName(cr) != "" {
		if key.CatalogName(cr) == defaultCatalogName {
			configMapName := fmt.Sprintf("%s-user-values", cr.Name)
			if key.UserConfigMapName(cr) != configMapName {
				return microerror.Maskf(validationError, "user configmap must be named %#q for app in default catalog", configMapName)
			}
		}

		ns := key.UserConfigMapNamespace(cr)
		if ns == "" {
			return microerror.Maskf(validationError, namespaceNotFoundReasonTemplate, "configmap", key.UserConfigMapName(cr))
		}

		_, err := v.k8sClient.CoreV1().ConfigMaps(ns).Get(ctx, key.UserConfigMapName(cr), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			return microerror.Maskf(validationError, resourceNotFoundTemplate, "configmap", key.UserConfigMapName(cr), ns)
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	if key.UserSecretName(cr) != "" {
		if key.CatalogName(cr) == defaultCatalogName {
			secretName := fmt.Sprintf("%s-user-secrets", cr.Name)
			if key.UserSecretName(cr) != secretName {
				return microerror.Maskf(validationError, "user secret must be named %#q for app in default catalog", secretName)
			}
		}

		ns := key.UserSecretNamespace(cr)
		if ns == "" {
			return microerror.Maskf(validationError, namespaceNotFoundReasonTemplate, "secret", key.UserSecretName(cr))
		}

		_, err := v.k8sClient.CoreV1().Secrets(key.UserSecretNamespace(cr)).Get(ctx, key.UserSecretName(cr), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			return microerror.Maskf(validationError, resourceNotFoundTemplate, "secret", key.UserSecretName(cr), ns)
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func contains(s []v1alpha1.Provider, e v1alpha1.Provider) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
