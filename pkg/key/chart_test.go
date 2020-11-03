package key

import (
	"testing"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_ChartStatus(t *testing.T) {
	expectedStatus := v1alpha1.ChartStatus{
		AppVersion: "0.12.0",
		Release: v1alpha1.ChartStatusRelease{
			Status: "DEPLOYED",
		},
		Version: "0.1.0",
	}

	obj := v1alpha1.Chart{
		Status: v1alpha1.ChartStatus{
			AppVersion: "0.12.0",
			Release: v1alpha1.ChartStatusRelease{
				Status: "DEPLOYED",
			},
			Version: "0.1.0",
		},
	}

	if ChartStatus(obj) != expectedStatus {
		t.Fatalf("chart status %#v, want %#v", ChartStatus(obj), expectedStatus)
	}
}

func Test_ChartConfigMapName(t *testing.T) {
	expectedName := "my-test-app-chart-values"

	obj := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-test-app",
			Namespace: "giantswarn",
		},
		Spec: v1alpha1.AppSpec{
			Name:    "test-app",
			Catalog: "test-catalog",
			Config: v1alpha1.AppSpecConfig{
				ConfigMap: v1alpha1.AppSpecConfigConfigMap{
					Name: "test-app-value",
				},
			},
		},
	}

	if ChartConfigMapName(obj) != expectedName {
		t.Fatalf("chartConfigMapName %#q, want %#q", ChartConfigMapName(obj), expectedName)
	}
}
