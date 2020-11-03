package key

import (
	"fmt"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"

	"github.com/giantswarm/app/v3/pkg/annotation"
)

func ChartConfigMapName(customResource v1alpha1.App) string {
	return fmt.Sprintf("%s-chart-values", customResource.GetName())
}

func ChartSecretName(customResource v1alpha1.App) string {
	return fmt.Sprintf("%s-chart-secrets", customResource.GetName())
}

func ChartStatus(customResource v1alpha1.Chart) v1alpha1.ChartStatus {
	return customResource.Status
}

func IsChartCordoned(customResource v1alpha1.Chart) bool {
	_, reasonOk := customResource.Annotations[fmt.Sprintf("%s/%s", annotation.ChartOperatorPrefix, annotation.CordonReason)]
	_, untilOk := customResource.Annotations[fmt.Sprintf("%s/%s", annotation.ChartOperatorPrefix, annotation.CordonUntil)]

	if reasonOk && untilOk {
		return true
	}

	return false
}
