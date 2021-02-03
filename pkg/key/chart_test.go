package key

import (
	"reflect"
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

func Test_ToChart(t *testing.T) {
	testCases := []struct {
		name           string
		input          interface{}
		expectedObject v1alpha1.Chart
		errorMatcher   func(error) bool
	}{
		{
			name: "case 0: basic match",
			input: &v1alpha1.Chart{
				Spec: v1alpha1.ChartSpec{
					Name:       "api",
					Namespace:  "giantswarm",
					TarballURL: "https://giantswarm.github.io/control-plane-catalog/api-1.0.0.tgz",
					Version:    "1.0.0",
				},
			},
			expectedObject: v1alpha1.Chart{
				Spec: v1alpha1.ChartSpec{
					Name:       "api",
					Namespace:  "giantswarm",
					TarballURL: "https://giantswarm.github.io/control-plane-catalog/api-1.0.0.tgz",
					Version:    "1.0.0",
				},
			},
		},
		{
			name:         "case 1: empty value",
			input:        nil,
			errorMatcher: IsEmptyValueError,
		},
		{
			name:         "case 2: wrong type",
			input:        &v1alpha1.AppCatalog{},
			errorMatcher: IsWrongTypeError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ToChart(tc.input)
			switch {
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case err != nil && !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if !reflect.DeepEqual(result, tc.expectedObject) {
				t.Fatalf("Custom Object == %#v, want %#v", result, tc.expectedObject)
			}
		})
	}
}
