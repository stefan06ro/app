package values

import (
	"context"
	"fmt"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/imdario/mergo"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/app/v4/pkg/key"
)

// MergeSecretData merges the data from the catalog, app and user secretss
// and returns a single set of values.
func (v *Values) MergeSecretData(ctx context.Context, app v1alpha1.App, appCatalog v1alpha1.AppCatalog) (map[string]interface{}, error) {
	appSecretName := key.AppSecretName(app)
	catalogSecretName := key.AppCatalogSecretName(appCatalog)
	userSecretName := key.UserSecretName(app)

	if appSecretName == "" && catalogSecretName == "" && userSecretName == "" {
		// Return early as there is no secret.
		return nil, nil
	}

	// We get the catalog level secrets if configured.
	rawCatalogData, err := v.getSecretDataForCatalog(ctx, appCatalog)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	catalogData, err := extractData(secret, "catalog", toStringMap(rawCatalogData))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// We get the app level secrets if configured.
	rawAppData, err := v.getSecretDataForApp(ctx, app)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	appData, err := extractData(secret, "app", toStringMap(rawAppData))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Secrets are merged and in case of intersecting values the app level
	// secrets are preferred.
	err = mergo.Merge(&catalogData, appData, mergo.WithOverride)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// We get the user level values if configured and merge them.
	if userSecretName != "" {
		rawUserData, err := v.getUserSecretDataForApp(ctx, app)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		// Secrets are merged again and in case of intersecting values the user
		// level secrets are preferred.
		userData, err := extractData(secret, "user", toStringMap(rawUserData))
		if err != nil {
			return nil, microerror.Mask(err)
		}

		err = mergo.Merge(&catalogData, userData, mergo.WithOverride)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return catalogData, nil
}

func (v *Values) getSecret(ctx context.Context, secretName, secretNamespace string) (map[string][]byte, error) {
	if secretName == "" {
		// Return early as no secret has been specified.
		return nil, nil
	}

	v.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("looking for secret %#q in namespace %#q", secretName, secretNamespace))

	secret, err := v.k8sClient.CoreV1().Secrets(secretNamespace).Get(ctx, secretName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil, microerror.Maskf(notFoundError, "secret %#q in namespace %#q not found", secretName, secretNamespace)
	} else if err != nil {
		return nil, microerror.Mask(err)
	}

	v.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found secret %#q in namespace %#q", secretName, secretNamespace))

	return secret.Data, nil
}

func (v *Values) getSecretDataForApp(ctx context.Context, app v1alpha1.App) (map[string][]byte, error) {
	secret, err := v.getSecret(ctx, key.AppSecretName(app), key.AppSecretNamespace(app))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return secret, nil
}

func (v *Values) getSecretDataForCatalog(ctx context.Context, catalog v1alpha1.AppCatalog) (map[string][]byte, error) {
	secret, err := v.getSecret(ctx, key.AppCatalogSecretName(catalog), key.AppCatalogSecretNamespace(catalog))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return secret, nil
}

func (v *Values) getUserSecretDataForApp(ctx context.Context, app v1alpha1.App) (map[string][]byte, error) {
	secret, err := v.getSecret(ctx, key.UserSecretName(app), key.UserSecretNamespace(app))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return secret, nil
}
