module mynet

go 1.16

require (
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/etcd v3.3.17+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gogo/protobuf v1.3.2
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.4.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/btree v1.0.1 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/raft v1.2.0
	github.com/hashicorp/raft-boltdb v0.0.0-20191021154308-4207f1bf0617
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/json-iterator/go v1.1.7
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nats-io/nats-server/v2 v2.6.2 // indirect
	github.com/nats-io/nats.go v1.13.0
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/etcd v3.3.17+incompatible
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	google.golang.org/protobuf v1.23.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
