module github.com/polarbroadband/rp1/dns_server

replace (
	github.com/polarbroadband/rp1/etcdlib => ../etcdlib
	github.com/polarbroadband/rp1/gitlib => ../gitlib
	github.com/polarbroadband/rp1/proto => ../proto
)

go 1.19

require (
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	go.etcd.io/etcd/api/v3 v3.5.6 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.6 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/tools v0.4.0 // indirect
	google.golang.org/genproto v0.0.0-20221227171554-f9683d7f8bef // indirect
	google.golang.org/grpc v1.51.0 // indirect
)

require (
	github.com/miekg/dns v1.1.50
	github.com/polarbroadband/rp1/etcdlib v0.0.0-00010101000000-000000000000
	github.com/polarbroadband/rp1/proto v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.0
	go.etcd.io/etcd/client/v3 v3.5.6
	google.golang.org/protobuf v1.28.1
)
