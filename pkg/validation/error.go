package validation

import (
	"regexp"
	"strings"

	"github.com/giantswarm/microerror"
)

const (
	appAdmissionControllerText string = "admission webhook \"apps.app-admission-controller-unique.giantswarm.io\" denied the request:"
)

var (
	appConfigMapNotFoundPattern = regexp.MustCompile(`app config map not found error: configmap [\d\D]+ in namespace [\d\D]+ not found`)
	kubeConfigNotFoundPattern   = regexp.MustCompile(`kube config not found error: kubeconfig secret [\d\D]+ in namespace [\d\D]+ not found`)
)

var appConfigMapNotFoundError = &microerror.Error{
	Kind: "appConfigMapNotFoundError",
}

// IsAppConfigMapNotFound asserts appConfigMapNotFoundError.
func IsAppConfigMapNotFound(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if strings.Contains(c.Error(), appAdmissionControllerText) && appConfigMapNotFoundPattern.MatchString(c.Error()) {
		return true
	}

	if c == appConfigMapNotFoundError { //nolint:gosimple
		return true
	}

	return false
}

var appDependencyNotReadyError = &microerror.Error{
	Kind: "appDependencyNotReadyError",
}

// IsAppDependencyNotReady asserts appDependencyNotReadyError.
func IsAppDependencyNotReady(err error) bool {
	return microerror.Cause(err) == appDependencyNotReadyError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var kubeConfigNotFoundError = &microerror.Error{
	Kind: "kubeConfigNotFoundError",
}

// IsKubeConfigNotFound asserts kubeConfigNotFoundError.
func IsKubeConfigNotFound(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if strings.Contains(c.Error(), appAdmissionControllerText) && kubeConfigNotFoundPattern.MatchString(c.Error()) {
		return true
	}

	if c == kubeConfigNotFoundError { //nolint:gosimple
		return true
	}

	return false
}

var notAllowedError = &microerror.Error{
	Kind: "notAllowedError",
}

// IsNotAllowed asserts notAllowedError.
func IsNotAllowed(err error) bool {
	return microerror.Cause(err) == notAllowedError
}

var notFoundError = &microerror.Error{
	Kind: "notFoundError",
}

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return microerror.Cause(err) == notFoundError
}

var parsingFailedError = &microerror.Error{
	Kind: "parsingFailedError",
}

// IsParsingFailed asserts parsingFailedError.
func IsParsingFailed(err error) bool {
	return microerror.Cause(err) == parsingFailedError
}

var validationError = &microerror.Error{
	Kind: "validationError",
}

// IsValidationError asserts validationError.
func IsValidationError(err error) bool {
	return microerror.Cause(err) == validationError
}
