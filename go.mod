module github.com/trackmyfish/backend

go 1.16

// Uncomment for local development
// replace github.com/trackmyfish/proto => ../proto

require (
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/jackc/pgproto3/v2 v2.1.0 // indirect
	github.com/jackc/pgx/v4 v4.11.0
	github.com/lib/pq v1.4.0 // indirect
	github.com/openlyinc/pointy v1.1.2
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v0.0.0-20200419222939-1884f454f8ea // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/trackmyfish/proto v0.0.9
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	google.golang.org/grpc v1.39.0
)
