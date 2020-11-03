package key

import (
	"testing"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
)

func Test_AppCatalogStorageURL(t *testing.T) {
	expectedURL := "http://giantswarm.io/sample-catalog/"

	obj := v1alpha1.AppCatalog{
		Spec: v1alpha1.AppCatalogSpec{
			Title:       "giant-swarm-title",
			Description: "giant-swarm app catalog sample",
			Storage: v1alpha1.AppCatalogSpecStorage{
				Type: "helm",
				URL:  "http://giantswarm.io/sample-catalog/",
			},
		},
	}

	if AppCatalogStorageURL(obj) != expectedURL {
		t.Fatalf("app catalog storage url %s, want %s", AppCatalogStorageURL(obj), expectedURL)
	}
}

func Test_AppCatalogTitle(t *testing.T) {
	expectedName := "giant-swarm-title"

	obj := v1alpha1.AppCatalog{
		Spec: v1alpha1.AppCatalogSpec{
			Title:       "giant-swarm-title",
			Description: "giant-swarm app catalog sample",
			Storage: v1alpha1.AppCatalogSpecStorage{
				Type: "helm",
				URL:  "http://giantswarm.io/sample-catalog.tgz",
			},
		},
	}

	if AppCatalogTitle(obj) != expectedName {
		t.Fatalf("app catalog name %s, want %s", AppCatalogTitle(obj), expectedName)
	}
}

func Test_AppCatalogConfigMapName(t *testing.T) {
	expectedName := "giant-swarm-configmap-name"

	obj := v1alpha1.AppCatalog{
		Spec: v1alpha1.AppCatalogSpec{
			Config: v1alpha1.AppCatalogSpecConfig{
				ConfigMap: v1alpha1.AppCatalogSpecConfigConfigMap{
					Name:      "giant-swarm-configmap-name",
					Namespace: "giant-swarm-configmap-namespace",
				},
			},
		},
	}

	if AppCatalogConfigMapName(obj) != expectedName {
		t.Fatalf("AppCatalogConfigMapName %#q, want %#q", AppCatalogConfigMapName(obj), expectedName)
	}
}

func Test_AppCatalogConfigMapNamespace(t *testing.T) {
	expectedNamespace := "giant-swarm-configmap-namespace"

	obj := v1alpha1.AppCatalog{
		Spec: v1alpha1.AppCatalogSpec{
			Config: v1alpha1.AppCatalogSpecConfig{
				ConfigMap: v1alpha1.AppCatalogSpecConfigConfigMap{
					Name:      "giant-swarm-configmap-name",
					Namespace: "giant-swarm-configmap-namespace",
				},
			},
		},
	}

	if AppCatalogConfigMapNamespace(obj) != expectedNamespace {
		t.Fatalf("AppCatalogConfigMapNamespace %#q, want %#q", AppCatalogConfigMapNamespace(obj), expectedNamespace)
	}
}

func Test_AppCatalogSecretName(t *testing.T) {
	expectedName := "giant-swarm-secret-name"

	obj := v1alpha1.AppCatalog{
		Spec: v1alpha1.AppCatalogSpec{
			Config: v1alpha1.AppCatalogSpecConfig{
				Secret: v1alpha1.AppCatalogSpecConfigSecret{
					Name:      "giant-swarm-secret-name",
					Namespace: "giant-swarm-secret-namespace",
				},
			},
		},
	}

	if AppCatalogSecretName(obj) != expectedName {
		t.Fatalf("AppCatalogSecretName %#q, want %#q", AppCatalogSecretName(obj), expectedName)
	}
}

func Test_AppCatalogSecretNamespace(t *testing.T) {
	expectedNamespace := "giant-swarm-secret-namespace"

	obj := v1alpha1.AppCatalog{
		Spec: v1alpha1.AppCatalogSpec{
			Config: v1alpha1.AppCatalogSpecConfig{
				Secret: v1alpha1.AppCatalogSpecConfigSecret{
					Name:      "giant-swarm-secret-name",
					Namespace: "giant-swarm-secret-namespace",
				},
			},
		},
	}

	if AppCatalogSecretNamespace(obj) != expectedNamespace {
		t.Fatalf("AppCatalogSecretNamespace %#q, want %#q", AppCatalogSecretNamespace(obj), expectedNamespace)
	}
}
