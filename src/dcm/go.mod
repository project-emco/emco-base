module gitlab.com/project-emco/core/emco-base/src/dcm

require (
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.18.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.0
	gitlab.com/project-emco/core/emco-base/src/monitor v0.0.0-00010101000000-000000000000
	gitlab.com/project-emco/core/emco-base/src/orchestrator v0.0.0-00010101000000-000000000000
	gitlab.com/project-emco/core/emco-base/src/rsync v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/contrib/propagators/b3 v1.8.0
	go.opentelemetry.io/otel v1.8.0
	go.opentelemetry.io/otel/exporters/zipkin v1.8.0
	go.opentelemetry.io/otel/sdk v1.8.0
	go.opentelemetry.io/otel/trace v1.8.0
	google.golang.org/grpc v1.43.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.23.3
	k8s.io/apimachinery v0.23.3
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	gitlab.com/project-emco/core/emco-base/src/clm => ../clm
	gitlab.com/project-emco/core/emco-base/src/monitor => ../monitor
	gitlab.com/project-emco/core/emco-base/src/orchestrator => ../orchestrator
	gitlab.com/project-emco/core/emco-base/src/rsync => ../rsync
	go.etcd.io/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200819165624-17cef6e3e9d5 // 17cef6e3e9d5 is the SHA for git tag v3.4.12
	google.golang.org/grpc => google.golang.org/grpc v1.29.0
)

go 1.16
