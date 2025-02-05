package key

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/annotation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_AppConfigMapName(t *testing.T) {
	testCases := []struct {
		name          string
		obj           v1alpha1.App
		expectedValue string
	}{
		{
			name:          "case 0: config is empty",
			obj:           v1alpha1.App{},
			expectedValue: "",
		},
		{
			name: "case 1: config has value",
			obj: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "configmap-name",
							Namespace: "configmap-namespace",
						},
					},
				},
			},
			expectedValue: "configmap-name",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			name := AppConfigMapName(tc.obj)

			if name != tc.expectedValue {
				t.Fatalf("AppConfigMapName %#q, want %#q", name, tc.expectedValue)
			}
		})
	}
}

func Test_AppConfigMapNamespace(t *testing.T) {
	testCases := []struct {
		name          string
		obj           v1alpha1.App
		expectedValue string
	}{
		{
			name:          "case 0: config is empty",
			obj:           v1alpha1.App{},
			expectedValue: "",
		},
		{
			name: "case 1: config has value",
			obj: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "configmap-name",
							Namespace: "configmap-namespace",
						},
					},
				},
			},
			expectedValue: "configmap-namespace",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			name := AppConfigMapNamespace(tc.obj)

			if name != tc.expectedValue {
				t.Fatalf("AppConfigMapNamespace %#q, want %#q", name, tc.expectedValue)
			}
		})
	}
}

func Test_AppName(t *testing.T) {
	expectedName := "giant-swarm-name"

	obj := v1alpha1.App{
		Spec: v1alpha1.AppSpec{
			Name: "giant-swarm-name",
		},
	}

	if AppName(obj) != expectedName {
		t.Fatalf("app name %#q, want %#q", AppName(obj), expectedName)
	}
}
func Test_AppSecretName(t *testing.T) {
	testCases := []struct {
		name          string
		obj           v1alpha1.App
		expectedValue string
	}{
		{
			name:          "case 0: config is empty",
			obj:           v1alpha1.App{},
			expectedValue: "",
		},
		{
			name: "case 1: config has value",
			obj: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Config: v1alpha1.AppSpecConfig{
						Secret: v1alpha1.AppSpecConfigSecret{
							Name:      "secret-name",
							Namespace: "secret-namespace",
						},
					},
				},
			},
			expectedValue: "secret-name",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			name := AppSecretName(tc.obj)

			if name != tc.expectedValue {
				t.Fatalf("AppSecretName %#q, want %#q", name, tc.expectedValue)
			}
		})
	}
}

func Test_AppSecretNamespace(t *testing.T) {
	testCases := []struct {
		name          string
		obj           v1alpha1.App
		expectedValue string
	}{
		{
			name:          "case 0: config is empty",
			obj:           v1alpha1.App{},
			expectedValue: "",
		},
		{
			name: "case 1: config has value",
			obj: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Config: v1alpha1.AppSpecConfig{
						Secret: v1alpha1.AppSpecConfigSecret{
							Name:      "secret-name",
							Namespace: "secret-namespace",
						},
					},
				},
			},
			expectedValue: "secret-namespace",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			name := AppSecretNamespace(tc.obj)

			if name != tc.expectedValue {
				t.Fatalf("AppSecretNamespace %#q, want %#q", name, tc.expectedValue)
			}
		})
	}
}

func Test_AppStatus(t *testing.T) {
	expectedStatus := v1alpha1.AppStatus{
		AppVersion: "0.12.0",
		Release: v1alpha1.AppStatusRelease{
			Status: "DEPLOYED",
		},
		Version: "0.1.0",
	}

	obj := v1alpha1.App{
		Status: v1alpha1.AppStatus{
			AppVersion: "0.12.0",
			Release: v1alpha1.AppStatusRelease{
				Status: "DEPLOYED",
			},
			Version: "0.1.0",
		},
	}

	if AppStatus(obj) != expectedStatus {
		t.Fatalf("app status %#q, want %#q", AppStatus(obj), expectedStatus)
	}
}

func Test_CatalogName(t *testing.T) {
	expectedName := "giant-swarm-catalog-name"

	obj := v1alpha1.App{
		Spec: v1alpha1.AppSpec{
			Name:    "giant-swarm-name",
			Catalog: "giant-swarm-catalog-name",
		},
	}

	if CatalogName(obj) != expectedName {
		t.Fatalf("catalog name %#q, want %#q", CatalogName(obj), expectedName)
	}
}

func Test_CordonReason(t *testing.T) {
	expectedCordonReason := "manual upgrade"

	obj := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				annotation.ChartOperatorCordonReason: "manual upgrade",
			},
		},
	}

	if CordonReason(obj) != expectedCordonReason {
		t.Fatalf("cordon reason %#q, want %s", CordonReason(obj), expectedCordonReason)
	}
}

func Test_CordonUntil(t *testing.T) {
	expectedCordonUntil := "2019-12-31T23:59:59Z"

	obj := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				annotation.ChartOperatorCordonUntil: "2019-12-31T23:59:59Z",
			},
		},
	}

	if CordonUntil(obj) != expectedCordonUntil {
		t.Fatalf("cordon until %s, want %s", CordonUntil(obj), expectedCordonUntil)
	}
}

func Test_InCluster(t *testing.T) {
	obj := v1alpha1.App{
		Spec: v1alpha1.AppSpec{
			KubeConfig: v1alpha1.AppSpecKubeConfig{
				InCluster: true,
			},
		},
	}

	if !InCluster(obj) {
		t.Fatalf("app namespace %#v, want %#v", InCluster(obj), true)
	}
}

func Test_IsAppCordoned(t *testing.T) {
	tests := []struct {
		name           string
		chart          v1alpha1.App
		expectedResult bool
	}{
		{
			name: "case 0: app cordoned",
			chart: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotation.AppOperatorCordonReason: "testing manual upgrade",
						annotation.AppOperatorCordonUntil:  "2019-12-31T23:59:59Z",
					},
				},
			},
			expectedResult: true,
		},
		{
			name:           "case 1: chart did not cordon",
			chart:          v1alpha1.App{},
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAppCordoned(tt.chart); got != tt.expectedResult {
				t.Errorf("IsCordoned() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}

func Test_KubeConfigFinalizer(t *testing.T) {
	obj := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-test-app",
		},
	}

	if KubeConfigFinalizer(obj) != "app-operator.giantswarm.io/app-my-test-app" {
		t.Fatalf("kubeconfig finalizer name %#v, want %#v", KubeConfigFinalizer(obj), "app-operator.giantswarm.io/app-my-test-app")
	}
}

func Test_KubeConfigSecretName(t *testing.T) {
	testCases := []struct {
		name          string
		obj           v1alpha1.App
		expectedValue string
	}{
		{
			name:          "case 0: config is empty",
			obj:           v1alpha1.App{},
			expectedValue: "",
		},
		{
			name: "case 1: config has value",
			obj: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "kubeconfig-name",
							Namespace: "kubeconfig-namespace",
						},
					},
				},
			},
			expectedValue: "kubeconfig-name",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			name := KubeConfigSecretName(tc.obj)

			if name != tc.expectedValue {
				t.Fatalf("KubeConfigSecretName %#q, want %#q", name, tc.expectedValue)
			}
		})
	}
}

func Test_KubeConfigSecretNamespace(t *testing.T) {
	testCases := []struct {
		name          string
		obj           v1alpha1.App
		expectedValue string
	}{
		{
			name:          "case 0: config is empty",
			obj:           v1alpha1.App{},
			expectedValue: "",
		},
		{
			name: "case 1: config has value",
			obj: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "kubeconfig-name",
							Namespace: "kubeconfig-namespace",
						},
					},
				},
			},
			expectedValue: "kubeconfig-namespace",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			name := KubeConfigSecretNamespace(tc.obj)

			if name != tc.expectedValue {
				t.Fatalf("KubeConfigSecretNamespace %#q, want %#q", name, tc.expectedValue)
			}
		})
	}
}

func Test_Namespace(t *testing.T) {
	expectedName := "giant-swarm-namespace"

	obj := v1alpha1.App{
		Spec: v1alpha1.AppSpec{
			Name:      "giant-swarm-name",
			Namespace: "giant-swarm-namespace",
		},
	}

	if Namespace(obj) != expectedName {
		t.Fatalf("app namespace %#q, want %#q", Namespace(obj), expectedName)
	}
}

func Test_ToApp(t *testing.T) {
	testCases := []struct {
		name           string
		input          interface{}
		expectedObject v1alpha1.App
		errorMatcher   func(error) bool
	}{
		{
			name: "case 0: basic match",
			input: &v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Name:      "giant-swarm-name",
					Namespace: "giant-swarm-namespace",
					Version:   "1.2.3",
				},
			},
			expectedObject: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Name:      "giant-swarm-name",
					Namespace: "giant-swarm-namespace",
					Version:   "1.2.3",
				},
			},
		},
		{
			name:         "case 1: wrong type",
			input:        &v1alpha1.AppCatalog{},
			errorMatcher: IsWrongTypeError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ToApp(tc.input)
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

func Test_UserConfigMapName(t *testing.T) {
	testCases := []struct {
		name          string
		obj           v1alpha1.App
		expectedValue string
	}{
		{
			name:          "case 0: config is empty",
			obj:           v1alpha1.App{},
			expectedValue: "",
		},
		{
			name: "case 1: config has value",
			obj: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "configmap-name",
							Namespace: "configmap-namespace",
						},
					},
				},
			},
			expectedValue: "configmap-name",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			name := UserConfigMapName(tc.obj)

			if name != tc.expectedValue {
				t.Fatalf("UserConfigMapName %#q, want %#q", name, tc.expectedValue)
			}
		})
	}
}

func Test_UserConfigMapNamespace(t *testing.T) {
	testCases := []struct {
		name          string
		obj           v1alpha1.App
		expectedValue string
	}{
		{
			name:          "case 0: config is empty",
			obj:           v1alpha1.App{},
			expectedValue: "",
		},
		{
			name: "case 1: config has value",
			obj: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "configmap-name",
							Namespace: "configmap-namespace",
						},
					},
				},
			},
			expectedValue: "configmap-namespace",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			name := UserConfigMapNamespace(tc.obj)

			if name != tc.expectedValue {
				t.Fatalf("UserConfigMapNamespace %#q, want %#q", name, tc.expectedValue)
			}
		})
	}
}

func Test_UserSecretName(t *testing.T) {
	testCases := []struct {
		name          string
		obj           v1alpha1.App
		expectedValue string
	}{
		{
			name:          "case 0: config is empty",
			obj:           v1alpha1.App{},
			expectedValue: "",
		},
		{
			name: "case 1: config has value",
			obj: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Name:    "name",
					Catalog: "catalog",
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "secret-name",
							Namespace: "secret-namespace",
						},
					},
				},
			},
			expectedValue: "secret-name",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			name := UserSecretName(tc.obj)

			if name != tc.expectedValue {
				t.Fatalf("UserSecretName %#q, want %#q", name, tc.expectedValue)
			}
		})
	}
}

func Test_UserSecretNamespace(t *testing.T) {
	testCases := []struct {
		name          string
		obj           v1alpha1.App
		expectedValue string
	}{
		{
			name:          "case 0: config is empty",
			obj:           v1alpha1.App{},
			expectedValue: "",
		},
		{
			name: "case 1: config has value",
			obj: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Name:    "name",
					Catalog: "catalog",
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "secret-name",
							Namespace: "secret-namespace",
						},
					},
				},
			},
			expectedValue: "secret-namespace",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			name := UserSecretNamespace(tc.obj)

			if name != tc.expectedValue {
				t.Fatalf("UserSecretNamespace %#q, want %#q", name, tc.expectedValue)
			}
		})
	}
}

func Test_Version(t *testing.T) {
	expectedVersion := "1.2.3"

	obj := v1alpha1.App{
		Spec: v1alpha1.AppSpec{
			Name:      "prometheus",
			Namespace: "monitoring",
			Version:   "1.2.3",
		},
	}

	if Version(obj) != expectedVersion {
		t.Fatalf("app version %#q, want %#q", Version(obj), expectedVersion)
	}
}

func Test_VersionLabel(t *testing.T) {
	testCases := []struct {
		name            string
		obj             v1alpha1.App
		expectedVersion string
		errorMatcher    func(error) bool
	}{
		{
			name: "case 0: basic match",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app-operator.giantswarm.io/version": "1.0.0",
					},
				},
			},
			expectedVersion: "1.0.0",
		},
		{
			name: "case 1: different value",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app-operator.giantswarm.io/version": "2.0.0",
					},
				},
			},
			expectedVersion: "2.0.0",
		},
		{
			name: "case 2: incorrect label",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"chart-operator.giantswarm.io/version": "1.0.0",
					},
				},
			},
			expectedVersion: "",
		},
		{
			name:            "case 3: no labels",
			obj:             v1alpha1.App{},
			expectedVersion: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := VersionLabel(tc.obj)

			if !reflect.DeepEqual(result, tc.expectedVersion) {
				t.Fatalf("Version label == %#v, want %#v", result, tc.expectedVersion)
			}
		})
	}
}
