package key

import (
	"fmt"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/annotation"
	"github.com/giantswarm/microerror"
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
	_, reasonOk := customResource.Annotations[annotation.ChartOperatorCordonReason]
	_, untilOk := customResource.Annotations[annotation.ChartOperatorCordonUntil]

	if reasonOk && untilOk {
		return true
	}

	return false
}

func ToChart(v interface{}) (v1alpha1.Chart, error) {
	if v == nil {
		return v1alpha1.Chart{}, microerror.Maskf(emptyValueError, "empty value cannot be converted to customResource")
	}

	customResource, ok := v.(*v1alpha1.Chart)
	if !ok {
		return v1alpha1.Chart{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha1.Chart{}, v)
	}

	return *customResource, nil
}
