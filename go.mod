module github.com/giantswarm/app/v5

go 1.16

require (
	github.com/giantswarm/apiextensions/v3 v3.32.0
	github.com/giantswarm/k8smetadata v0.3.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/micrologger v0.5.0
	github.com/giantswarm/to v0.3.0
	github.com/google/go-cmp v0.5.6
	github.com/google/go-github/v35 v35.3.0
	github.com/imdario/mergo v0.3.12
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	k8s.io/api v0.20.11
	k8s.io/apiextensions-apiserver v0.20.11
	k8s.io/apimachinery v0.20.11
	k8s.io/client-go v0.20.11
	sigs.k8s.io/yaml v1.3.0
)

replace (
	github.com/Microsoft/hcsshim v0.8.7 => github.com/Microsoft/hcsshim v0.8.10
	github.com/bketelsen/crypt => github.com/bketelsen/crypt v0.0.3
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	// Use moby v20.10.0-beta1 to fix build issue on darwin.
	github.com/docker/docker => github.com/moby/moby v20.10.0-beta1+incompatible
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc7
	github.com/spf13/viper => github.com/spf13/viper v1.7.1
	// Use fork of CAPI with Kubernetes 1.18 support.
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
)
