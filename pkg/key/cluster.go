package key

import (
	"fmt"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
)

const (
	ingressControllerConfigMapName = "ingress-controller-values"
	nginxIngressControllerAppName  = "nginx-ingress-controller-app"
)

func ClusterConfigMapName(customResource v1alpha1.App) string {
	// A separate config map is used for Nginx Ingress Controller.
	if AppName(customResource) == nginxIngressControllerAppName {
		return ingressControllerConfigMapName
	}

	return fmt.Sprintf("%s-cluster-values", customResource.Namespace)
}

func ClusterKubeConfigSecretName(customResource v1alpha1.App) string {
	return fmt.Sprintf("%s-kubeconfig", customResource.Namespace)
}
