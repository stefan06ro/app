package validation

import (
	"context"
	"strings"
	"testing"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/apiextensions/v3/pkg/clientset/versioned/fake"
	"github.com/giantswarm/apiextensions/v3/pkg/label"
	"github.com/giantswarm/micrologger/microloggertest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgofake "k8s.io/client-go/kubernetes/fake"
)

func Test_ValidateApp(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		obj         v1alpha1.App
		catalogs    []*v1alpha1.AppCatalog
		configMaps  []*corev1.ConfigMap
		secrets     []*corev1.Secret
		expectedErr string
	}{
		{
			name: "case 0: flawless flow",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "eggs2-cluster-values",
							Namespace: "eggs2",
						},
					},
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						Context: v1alpha1.AppSpecKubeConfigContext{
							Name: "eggs2-kubeconfig",
						},
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "eggs2-kubeconfig",
							Namespace: "eggs2",
						},
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "kiam-user-values",
							Namespace: "eggs2",
						},
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("giantswarm"),
			},
			configMaps: []*corev1.ConfigMap{
				newTestConfigMap("eggs2-cluster-values", "eggs2"),
				newTestConfigMap("kiam-user-values", "eggs2"),
			},
			secrets: []*corev1.Secret{
				newTestSecret("eggs2-kubeconfig", "eggs2"),
			},
		},
		{
			name: "case 1: flawless in-cluster",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("control-plane-catalog"),
			},
		},
		{
			name: "case 2: missing version label",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("control-plane-catalog"),
			},
			expectedErr: "validation error: label `app-operator.giantswarm.io/version` not found",
		},
		{
			name: "case 3: spec.catalog not found",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			expectedErr: "validation error: catalog `control-plane-catalog` not found",
		},
		{
			name: "case 4: spec.config.configMap not found",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "dex-app-values",
							Namespace: "giantswarm",
						},
					},
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("control-plane-catalog"),
			},
			expectedErr: "app config map not found error: configmap `dex-app-values` in namespace `giantswarm` not found",
		},
		{
			name: "case 5: spec.config.configMap no namespace specified",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "dex-app-values",
							Namespace: "",
						},
					},
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("control-plane-catalog"),
			},
			expectedErr: "validation error: namespace is not specified for configmap `dex-app-values`",
		},
		{
			name: "case 6: spec.config.secret not found",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					Config: v1alpha1.AppSpecConfig{
						Secret: v1alpha1.AppSpecConfigSecret{
							Name:      "dex-app-secrets",
							Namespace: "giantswarm",
						},
					},
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("control-plane-catalog"),
			},
			expectedErr: "validation error: secret `dex-app-secrets` in namespace `giantswarm` not found",
		},
		{
			name: "case 7: spec.config.secret no namespace specified",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					Config: v1alpha1.AppSpecConfig{
						Secret: v1alpha1.AppSpecConfigSecret{
							Name:      "dex-app-secrets",
							Namespace: "",
						},
					},
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("control-plane-catalog"),
			},
			expectedErr: "validation error: namespace is not specified for secret `dex-app-secrets`",
		},
		{
			name: "case 8: spec.kubeConfig.secret not found",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						Context: v1alpha1.AppSpecKubeConfigContext{
							Name: "eggs2-kubeconfig",
						},
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "eggs2-kubeconfig",
							Namespace: "eggs2",
						},
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("giantswarm"),
			},
			expectedErr: "kube config not found error: kubeconfig secret `eggs2-kubeconfig` in namespace `eggs2` not found",
		},
		{
			name: "case 9: spec.kubeConfig.secret no namespace specified",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						Context: v1alpha1.AppSpecKubeConfigContext{
							Name: "eggs2-kubeconfig",
						},
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "eggs2-kubeconfig",
							Namespace: "",
						},
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("giantswarm"),
			},
			expectedErr: "validation error: namespace is not specified for kubeconfig secret `eggs2-kubeconfig`",
		},
		{
			name: "case 10: spec.userConfig.configMap not found",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "dex-app-user-values",
							Namespace: "giantswarm",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("control-plane-catalog"),
			},
			expectedErr: "validation error: configmap `dex-app-user-values` in namespace `giantswarm` not found",
		},
		{
			name: "case 11: spec.userConfig.configMap no namespace specified",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "dex-app-user-values",
							Namespace: "",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("control-plane-catalog"),
			},
			expectedErr: "validation error: namespace is not specified for configmap `dex-app-user-values`",
		},
		{
			name: "case 12: spec.userConfig.secret not found",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "dex-app-user-secrets",
							Namespace: "giantswarm",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("control-plane-catalog"),
			},
			expectedErr: "validation error: secret `dex-app-user-secrets` in namespace `giantswarm` not found",
		},
		{
			name: "case 13: spec.userConfig.secret no namespace specified",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "dex-app-user-secrets",
							Namespace: "",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.AppCatalog{
				newTestCatalog("control-plane-catalog"),
			},
			expectedErr: "validation error: namespace is not specified for secret `dex-app-user-secrets`",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g8sObjs := make([]runtime.Object, 0)
			for _, cat := range tc.catalogs {
				g8sObjs = append(g8sObjs, cat)
			}

			k8sObjs := make([]runtime.Object, 0)
			for _, cm := range tc.configMaps {
				k8sObjs = append(k8sObjs, cm)
			}

			for _, secret := range tc.secrets {
				k8sObjs = append(k8sObjs, secret)
			}

			c := Config{
				G8sClient: fake.NewSimpleClientset(g8sObjs...),
				K8sClient: clientgofake.NewSimpleClientset(k8sObjs...),
				Logger:    microloggertest.New(),

				Provider: "aws",
			}
			r, err := NewValidator(c)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			_, err = r.ValidateApp(ctx, tc.obj)
			switch {
			case err != nil && tc.expectedErr == "":
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.expectedErr != "":
				t.Fatalf("error == nil, want non-nil")
			}

			if err != nil && tc.expectedErr != "" {
				if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Fatalf("error == %#v, want %#v ", err.Error(), tc.expectedErr)
				}

			}
		})
	}
}

func Test_ValidateMetadataConstraints(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		obj          v1alpha1.App
		catalogEntry *v1alpha1.AppCatalogEntry
		apps         []*v1alpha1.App
		expectedErr  string
	}{
		{
			name: "case 0: flawless flow",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: metav1.NamespaceDefault,
					Version:   "1.4.0",
				},
			},
			catalogEntry: &v1alpha1.AppCatalogEntry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "giantswarm-kiam-1.4.0",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: v1alpha1.AppCatalogEntrySpec{
					Restrictions: &v1alpha1.AppCatalogEntrySpecRestrictions{
						FixedNamespace: metav1.NamespaceDefault,
					},
				},
			},
		},
		{
			name: "case 1: fixed namespace constraint",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Version:   "1.4.0",
				},
			},
			catalogEntry: &v1alpha1.AppCatalogEntry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "giantswarm-kiam-1.4.0",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: v1alpha1.AppCatalogEntrySpec{
					Restrictions: &v1alpha1.AppCatalogEntrySpecRestrictions{
						FixedNamespace: "eggs1",
					},
				},
			},
			expectedErr: "validation error: app `kiam` can only be installed in namespace `eggs1` only, not `kube-system`",
		},
		{
			name: "case 2: cluster singleton constraint",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Version:   "1.4.0",
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-kiam",
						Namespace: "eggs2",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "kiam",
						Namespace: "giantswarm",
						Version:   "1.3.0-rc1",
					},
				},
			},
			catalogEntry: &v1alpha1.AppCatalogEntry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "giantswarm-kiam-1.4.0",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: v1alpha1.AppCatalogEntrySpec{
					Restrictions: &v1alpha1.AppCatalogEntrySpecRestrictions{
						ClusterSingleton: true,
					},
				},
			},
			expectedErr: "validation error: app `kiam` can only be installed once in cluster `eggs2`",
		},
		{
			name: "case 3: namespace singleton constraint",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Version:   "1.4.0",
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-kiam",
						Namespace: "eggs2",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "kiam",
						Namespace: "giantswarm",
						Version:   "1.3.0-rc1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-kiam-1",
						Namespace: "eggs2",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "kiam",
						Namespace: "kube-system",
						Version:   "1.3.0-rc1",
					},
				},
			},
			catalogEntry: &v1alpha1.AppCatalogEntry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "giantswarm-kiam-1.4.0",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: v1alpha1.AppCatalogEntrySpec{
					Restrictions: &v1alpha1.AppCatalogEntrySpecRestrictions{
						NamespaceSingleton: true,
					},
				},
			},
			expectedErr: "validation error: app `kiam` can only be installed only once in namespace `kube-system`",
		},
		{
			name: "case 4: compatible providers constraint",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Version:   "1.4.0",
				},
			},
			catalogEntry: &v1alpha1.AppCatalogEntry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "giantswarm-kiam-1.4.0",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: v1alpha1.AppCatalogEntrySpec{
					Restrictions: &v1alpha1.AppCatalogEntrySpecRestrictions{
						CompatibleProviders: []v1alpha1.Provider{"azure"},
					},
				},
			},
			expectedErr: "validation error: app `kiam` can only be installed for providers [`azure`] not `aws`",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g8sObjs := make([]runtime.Object, 0)

			if tc.catalogEntry != nil {
				g8sObjs = append(g8sObjs, tc.catalogEntry)
			}

			for _, app := range tc.apps {
				g8sObjs = append(g8sObjs, app)
			}

			c := Config{
				G8sClient: fake.NewSimpleClientset(g8sObjs...),
				K8sClient: clientgofake.NewSimpleClientset(),
				Logger:    microloggertest.New(),

				Provider: "aws",
			}
			r, err := NewValidator(c)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			err = r.validateMetadataConstraints(ctx, tc.obj)
			switch {
			case err != nil && tc.expectedErr == "":
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.expectedErr != "":
				t.Fatalf("error == nil, want non-nil")
			}

			if err != nil && tc.expectedErr != "" {
				if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Fatalf("error == %#v, want %#v ", err.Error(), tc.expectedErr)
				}

			}
		})
	}
}

func newTestCatalog(name string) *v1alpha1.AppCatalog {
	return &v1alpha1.AppCatalog{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1alpha1.AppCatalogSpec{
			Description: name,
			Title:       name,
		},
	}
}

func newTestConfigMap(name, namespace string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		Data: map[string]string{
			"values": "cluster: yaml\n",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func newTestSecret(name, namespace string) *corev1.Secret {
	return &corev1.Secret{
		Data: map[string][]byte{
			"values": []byte("cluster: yaml\n"),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}
