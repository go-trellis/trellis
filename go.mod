module github.com/iTrellis/trellis

go 1.13

replace (
	go.etcd.io/etcd/api/v3 v3.5.0-pre => go.etcd.io/etcd/api/v3 v3.0.0-20210107172604-c632042bb96c
	go.etcd.io/etcd/pkg/v3 v3.5.0-pre => go.etcd.io/etcd/pkg/v3 v3.0.0-20210107172604-c632042bb96c
)

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/gzip v0.0.3
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.2 // indirect
	github.com/go-resty/resty/v2 v2.5.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.2.0
	github.com/iTrellis/common v0.21.7
	github.com/iTrellis/config v0.21.5
	github.com/iTrellis/node v0.21.4
	github.com/iTrellis/xorm_ext v0.21.5
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.1
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
	github.com/sirupsen/logrus v1.8.0
	github.com/ugorji/go v1.2.5 // indirect
	github.com/urfave/cli/v2 v2.3.0
	go.etcd.io/etcd/api/v3 v3.5.0-pre
	go.etcd.io/etcd/client/v3 v3.0.0-20210201223203-e897daaebc2f
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/sys v0.0.0-20210403161142-5e06dd20ab57 // indirect
	google.golang.org/grpc v1.29.1
	xorm.io/builder v0.3.9 // indirect
	xorm.io/xorm v1.0.7
)
