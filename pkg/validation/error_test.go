package validation

import (
	"errors"
	"testing"
)

func Test_IsAppConfigMapNotFound(t *testing.T) {
	tests := []struct {
		name          string
		errorMessage  string
		expectedMatch bool
	}{
		{
			name:          "case 0: good match",
			errorMessage:  "admission webhook \"apps.app-admission-controller-unique.giantswarm.io\" denied the request: app config map not found error: configmap `u9q0r-cluster-values` in namespace `u9q0r` not found",
			expectedMatch: true,
		},
		{
			name:          "case 1: wrong error message",
			errorMessage:  "admission webhook \"apps.app-admission-controller-unique.giantswarm.io\" denied the request: validation error: kubeconfig secret `u9q0r-cluster-values` in namespace `u9q0r` not found",
			expectedMatch: false,
		},
		{
			name:          "case 2: empty string",
			errorMessage:  "",
			expectedMatch: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := errors.New(tc.errorMessage)
			result := IsAppConfigMapNotFound(err)

			if result != tc.expectedMatch {
				t.Fatalf("expected %t, got %t", tc.expectedMatch, result)
			}
		})
	}
}

func Test_IsKubeConfigNotFound(t *testing.T) {
	tests := []struct {
		name          string
		errorMessage  string
		expectedMatch bool
	}{
		{
			name:          "case 0: good match",
			errorMessage:  "admission webhook \"apps.app-admission-controller-unique.giantswarm.io\" denied the request: kube config not found error: kubeconfig secret `u9q0r-kubeconfig` in namespace `u9q0r` not found",
			expectedMatch: true,
		},
		{
			name:          "case 1: wrong error message",
			errorMessage:  "admission webhook \"apps.app-admission-controller-unique.giantswarm.io\" denied the request: validation error: configmap `u9q0r-cluster-values` in namespace `u9q0r` not found",
			expectedMatch: false,
		},
		{
			name:          "case 2: empty string",
			errorMessage:  "",
			expectedMatch: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := errors.New(tc.errorMessage)
			result := IsKubeConfigNotFound(err)

			if result != tc.expectedMatch {
				t.Fatalf("expected %t, got %t", tc.expectedMatch, result)
			}
		})
	}
}
