module github.com/giantswarm/app/v3

go 1.15

require (
	github.com/giantswarm/apiextensions/v3 v3.6.0
	github.com/giantswarm/microerror v0.2.1
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	k8s.io/apimachinery v0.18.9
	sigs.k8s.io/yaml v1.2.0
)

replace (
	// Apply fix for CVE-2020-15114 not yet released in github.com/spf13/viper.
	github.com/bketelsen/crypt => github.com/bketelsen/crypt v0.0.3
	// Apply security fix not present in v1.4.0.
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	// Apply security fix not present in 1.6.2.
	github.com/spf13/viper => github.com/spf13/viper v1.7.1
	// Use fork of CAPI with Kubernetes 1.18 support.
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
)
