package key

import (
	"fmt"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"

	"github.com/giantswarm/app/v4/pkg/annotation"
)

func AppCatalogEntryManagedBy(projectName string) string {
	return fmt.Sprintf("%s-unique", projectName)
}

func AppCatalogEntryName(catalogName, appName, appVersion string) string {
	return fmt.Sprintf("%s-%s-%s", catalogName, appName, appVersion)
}

func AppCatalogEntryTeam(customResource v1alpha1.AppCatalogEntry) string {
	return customResource.Annotations[annotation.Team]
}
