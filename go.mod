module kubedb.dev/cli

go 1.23.0

toolchain go1.24.1

require (
	github.com/Masterminds/semver/v3 v3.3.1
	github.com/cert-manager/cert-manager v1.17.1
	github.com/fatih/camelcase v1.0.0
	github.com/go-sql-driver/mysql v1.9.0
	github.com/prometheus/client_golang v1.20.5
	github.com/prometheus/common v0.61.0
	github.com/spf13/cobra v1.9.1
	golang.org/x/text v0.23.0
	gomodules.xyz/go-sh v0.2.0
	gomodules.xyz/logs v0.0.7
	gomodules.xyz/pointer v0.1.0
	gomodules.xyz/runtime v0.3.0
	gomodules.xyz/x v0.0.17
	k8s.io/api v0.32.3
	k8s.io/apimachinery v0.32.3
	k8s.io/cli-runtime v0.32.2
	k8s.io/client-go v0.32.3
	k8s.io/component-base v0.32.3
	k8s.io/klog/v2 v2.130.1
	k8s.io/kubectl v0.29.0
	kmodules.xyz/cert-manager-util v0.29.0
	kmodules.xyz/client-go v0.32.1
	kmodules.xyz/custom-resources v0.32.0
	kmodules.xyz/monitoring-agent-api v0.32.0
	kubedb.dev/apimachinery v0.53.0-rc.1
	kubedb.dev/db-client-go v0.8.0-rc.1
	kubeops.dev/petset v0.0.9
	sigs.k8s.io/controller-runtime v0.20.3
	sigs.k8s.io/yaml v1.4.0
	stash.appscode.dev/apimachinery v0.39.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chai2010/gettext-go v1.0.2 // indirect
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/cyphar/filepath-securejoin v0.3.4 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elastic/elastic-transport-go/v8 v8.1.0 // indirect
	github.com/elastic/go-elasticsearch/v5 v5.6.1 // indirect
	github.com/elastic/go-elasticsearch/v6 v6.8.10 // indirect
	github.com/elastic/go-elasticsearch/v7 v7.15.1 // indirect
	github.com/elastic/go-elasticsearch/v8 v8.4.0 // indirect
	github.com/emicklei/go-restful/v3 v3.12.1 // indirect
	github.com/evanphx/json-patch v5.9.11+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-resty/resty/v2 v2.11.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/gnostic-models v0.6.9 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-containerregistry v0.20.3 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/kubernetes-csi/external-snapshotter/client/v7 v7.0.0 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/spdystream v0.5.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/onsi/gomega v1.36.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opensearch-project/opensearch-go v1.1.0 // indirect
	github.com/opensearch-project/opensearch-go/v2 v2.3.0 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.81.0 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.81.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/redis/go-redis/v9 v9.5.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xlab/treeprint v1.2.0 // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.virtual-secrets.dev/apimachinery v0.0.1 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/oauth2 v0.27.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/time v0.10.0 // indirect
	gomodules.xyz/clock v0.0.0-20200817085942-06523dba733f // indirect
	gomodules.xyz/encoding v0.0.8 // indirect
	gomodules.xyz/flags v0.1.3 // indirect
	gomodules.xyz/jsonpatch/v2 v2.5.0 // indirect
	gomodules.xyz/mergo v0.3.13 // indirect
	gomodules.xyz/password-generator v0.2.9 // indirect
	gomodules.xyz/sets v0.2.1 // indirect
	gomodules.xyz/wait v0.2.0 // indirect
	google.golang.org/protobuf v1.36.3 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.32.3 // indirect
	k8s.io/apiserver v0.32.3 // indirect
	k8s.io/kube-aggregator v0.32.3 // indirect
	k8s.io/kube-openapi v0.0.0-20250318172550-b98be4ee1595 // indirect
	k8s.io/metrics v0.32.3 // indirect
	k8s.io/utils v0.0.0-20241210054802-24370beab758 // indirect
	kmodules.xyz/apiversion v0.2.0 // indirect
	kmodules.xyz/objectstore-api v0.32.0 // indirect
	kmodules.xyz/offshoot-api v0.32.0 // indirect
	kmodules.xyz/prober v0.32.0 // indirect
	kmodules.xyz/resource-metadata v0.26.1 // indirect
	kubeops.dev/csi-driver-cacerts v0.1.0 // indirect
	kubeops.dev/sidekick v0.0.10 // indirect
	kubestash.dev/apimachinery v0.17.0-rc.0 // indirect
	sigs.k8s.io/gateway-api v1.1.0 // indirect
	sigs.k8s.io/json v0.0.0-20241014173422-cfa47c3a1cc8 // indirect
	sigs.k8s.io/kustomize/api v0.18.0 // indirect
	sigs.k8s.io/kustomize/kyaml v0.18.1 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.6.0 // indirect
)

replace github.com/Masterminds/sprig/v3 => github.com/gomodules/sprig/v3 v3.2.3-0.20220405051441-0a8a99bac1b8

replace sigs.k8s.io/controller-runtime => github.com/kmodules/controller-runtime v0.20.3-0.20250221050548-8eabe54e7dda

replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.6

replace k8s.io/apiserver => github.com/kmodules/apiserver v0.32.3-0.20250221062720-35dc674c7dd6
