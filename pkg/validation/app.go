package validation

import (
	"context"

	"github.com/giantswarm/apiextensions/v3/pkg/annotation"
	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/apiextensions/v3/pkg/label"
	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/app/v4/pkg/key"
)

const (
	catalogNotFoundTemplate         = "catalog %#q not found"
	labelNotFoundTemplate           = "label %#q not found"
	namespaceNotFoundReasonTemplate = "namespace is not specified for %s %#q"
	resourceNotFoundTemplate        = "%s %#q in namespace %#q not found"
)

func (v *Validator) ValidateApp(ctx context.Context, app v1alpha1.App) (bool, error) {
	var err error

	err = v.validateLabels(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

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

	err = v.validateUserConfig(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	return true, nil
}

func (v *Validator) validateCatalog(ctx context.Context, cr v1alpha1.App) error {
	if key.CatalogName(cr) != "" {
		_, err := v.g8sClient.ApplicationV1alpha1().AppCatalogs().Get(ctx, key.CatalogName(cr), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			return microerror.Maskf(validationError, catalogNotFoundTemplate, key.CatalogName(cr))
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (v *Validator) validateConfig(ctx context.Context, cr v1alpha1.App) error {
	_, hasManagedConfig := cr.Annotations[annotation.ConfigVersion]
	if hasManagedConfig && (key.AppConfigMapName(cr) == "" || key.AppSecretName(cr) == "") {
		// wait for config-controller setting app CR configmap and secret
		return microerror.Maskf(appDependencyNotReadyError, "ConfigMap or Secret not set")
	}

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

func (v *Validator) validateLabels(ctx context.Context, cr v1alpha1.App) error {
	if key.VersionLabel(cr) == "" {
		return microerror.Maskf(validationError, labelNotFoundTemplate, label.AppOperatorVersion)
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

func (v *Validator) validateUserConfig(ctx context.Context, cr v1alpha1.App) error {
	if key.UserConfigMapName(cr) != "" {
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
