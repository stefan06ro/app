package key

import (
	"testing"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
)

func Test_CatalogStorageURL(t *testing.T) {
	expectedURL := "http://giantswarm.io/sample-catalog/"

	obj := v1alpha1.Catalog{
		Spec: v1alpha1.CatalogSpec{
			Title:       "giant-swarm-title",
			Description: "giant-swarm app catalog sample",
			Storage: v1alpha1.CatalogSpecStorage{
				Type: "helm",
				URL:  "http://giantswarm.io/sample-catalog/",
			},
		},
	}

	if CatalogStorageURL(obj) != expectedURL {
		t.Fatalf("app catalog storage url %s, want %s", CatalogStorageURL(obj), expectedURL)
	}
}

func Test_CatalogTitle(t *testing.T) {
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

func Test_CatalogConfigMapName(t *testing.T) {
	expectedName := "giant-swarm-configmap-name"

	obj := v1alpha1.Catalog{
		Spec: v1alpha1.CatalogSpec{
			Config: &v1alpha1.CatalogSpecConfig{
				ConfigMap: &v1alpha1.CatalogSpecConfigConfigMap{
					Name:      "giant-swarm-configmap-name",
					Namespace: "giant-swarm-configmap-namespace",
				},
			},
		},
	}

	if CatalogConfigMapName(obj) != expectedName {
		t.Fatalf("CatalogConfigMapName %#q, want %#q", CatalogConfigMapName(obj), expectedName)
	}
}

func Test_CatalogConfigMapNamespace(t *testing.T) {
	expectedNamespace := "giant-swarm-configmap-namespace"

	obj := v1alpha1.Catalog{
		Spec: v1alpha1.CatalogSpec{
			Config: &v1alpha1.CatalogSpecConfig{
				ConfigMap: &v1alpha1.CatalogSpecConfigConfigMap{
					Name:      "giant-swarm-configmap-name",
					Namespace: "giant-swarm-configmap-namespace",
				},
			},
		},
	}

	if CatalogConfigMapNamespace(obj) != expectedNamespace {
		t.Fatalf("CatalogConfigMapNamespace %#q, want %#q", CatalogConfigMapNamespace(obj), expectedNamespace)
	}
}

func Test_CatalogSecretName(t *testing.T) {
	expectedName := "giant-swarm-secret-name"

	obj := v1alpha1.Catalog{
		Spec: v1alpha1.CatalogSpec{
			Config: &v1alpha1.CatalogSpecConfig{
				Secret: &v1alpha1.CatalogSpecConfigSecret{
					Name:      "giant-swarm-secret-name",
					Namespace: "giant-swarm-secret-namespace",
				},
			},
		},
	}

	if CatalogSecretName(obj) != expectedName {
		t.Fatalf("CatalogSecretName %#q, want %#q", CatalogSecretName(obj), expectedName)
	}
}

func Test_CatalogSecretNamespace(t *testing.T) {
	expectedNamespace := "giant-swarm-secret-namespace"

	obj := v1alpha1.Catalog{
		Spec: v1alpha1.CatalogSpec{
			Config: &v1alpha1.CatalogSpecConfig{
				Secret: &v1alpha1.CatalogSpecConfigSecret{
					Name:      "giant-swarm-secret-name",
					Namespace: "giant-swarm-secret-namespace",
				},
			},
		},
	}

	if CatalogSecretNamespace(obj) != expectedNamespace {
		t.Fatalf("CatalogSecretNamespace %#q, want %#q", CatalogSecretNamespace(obj), expectedNamespace)
	}
}
