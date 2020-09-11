module github.com/argoproj/argo-rollouts

go 1.13

require (
	github.com/Djarvur/go-err113 v0.1.0 // indirect
	github.com/antonmedv/expr v1.4.2
	github.com/argoproj/pkg v0.0.0-20200624215116-23e74cb168fe
	github.com/bouk/monkey v1.0.0
	github.com/cucumber/godog v0.10.0
	github.com/daixiang0/gci v0.2.3 // indirect
	github.com/docker/docker v1.4.2-0.20190327010347-be7ac8be2ae0 // indirect
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-critic/go-critic v0.5.2 // indirect
	github.com/go-openapi/spec v0.19.3
	github.com/gofrs/flock v0.8.0 // indirect
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/golang/mock v1.4.4 // indirect
	github.com/golangci/golangci-lint v1.30.0 // indirect
	github.com/golangci/misspell v0.3.5 // indirect
	github.com/golangci/revgrep v0.0.0-20180812185044-276a5c0a1039 // indirect
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/gostaticanalysis/analysisutil v0.1.0 // indirect
	github.com/gostaticanalysis/comment v1.4.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/jirfag/go-printf-func-name v0.0.0-20200119135958-7558a9eaa5af // indirect
	github.com/jstemmer/go-junit-report v0.9.1
	github.com/juju/ansiterm v0.0.0-20180109212912-720a0952cc2a
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/magiconair/properties v1.8.2 // indirect
	github.com/matoous/godox v0.0.0-20200801072554-4fb83dc2941e // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/nishanths/exhaustive v0.0.0-20200811152831-6cf413ae40e0 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.0
	github.com/prometheus/common v0.9.1
	github.com/quasilyte/regex/syntax v0.0.0-20200805063351-8f842688393c // indirect
	github.com/servicemeshinterface/smi-sdk-go v0.3.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spaceapegames/go-wavefront v1.8.1
	github.com/spf13/afero v1.3.4 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1 // indirect
	github.com/ssgreg/nlreturn/v2 v2.1.0 // indirect
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/tdakkota/asciicheck v0.0.0-20200416200610-e657995f937b // indirect
	github.com/timakin/bodyclose v0.0.0-20200424151742-cb6215831a94 // indirect
	github.com/ultraware/funlen v0.0.3 // indirect
	github.com/valyala/fasttemplate v1.2.1
	github.com/vektra/mockery v1.1.2
	golang.org/x/sys v0.0.0-20200828194041-157a740278f4 // indirect
	golang.org/x/tools v0.0.0-20200828161849-5deb26317202 // indirect
	gopkg.in/ini.v1 v1.60.2 // indirect
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	gotest.tools v2.2.0+incompatible
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
	k8s.io/api v0.17.4
	k8s.io/apiextensions-apiserver v0.17.0
	k8s.io/apimachinery v0.17.4
	k8s.io/apiserver v0.17.3
	k8s.io/cli-runtime v0.17.3
	k8s.io/client-go v0.17.4
	k8s.io/code-generator v0.17.4
	k8s.io/component-base v0.17.3
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a
	k8s.io/kubectl v0.16.4
	k8s.io/kubernetes v1.17.3
	k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
	mvdan.cc/gofumpt v0.0.0-20200802201014-ab5a8192947d // indirect
	mvdan.cc/unparam v0.0.0-20200501210554-b37ab49443f7 // indirect
	sigs.k8s.io/controller-tools v0.2.5
	sourcegraph.com/sqs/pbtypes v1.0.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.17.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.4-beta.0
	k8s.io/apiserver => k8s.io/apiserver v0.17.3
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.3
	k8s.io/client-go => k8s.io/client-go v0.17.3
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.3
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.3
	k8s.io/code-generator => k8s.io/code-generator v0.17.4-beta.0
	k8s.io/component-base => k8s.io/component-base v0.17.3
	k8s.io/cri-api => k8s.io/cri-api v0.17.4-beta.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.3
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.17.3
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.3
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.3
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.3
	k8s.io/kubectl => k8s.io/kubectl v0.17.3
	k8s.io/kubelet => k8s.io/kubelet v0.17.3
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.17.3
	k8s.io/metrics => k8s.io/metrics v0.17.3
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.17.3
)
