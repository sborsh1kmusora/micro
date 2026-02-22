module github.com/sborsh1kmusora/micro/inventory

go 1.25.3

require (
	github.com/google/uuid v1.6.0
	github.com/sborsh1kmusora/micro/shared v0.0.0-20260207124816-dd99df575cd0
	google.golang.org/grpc v1.78.0
)

require (
	github.com/brianvoe/gofakeit/v7 v7.14.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.7 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260128011058-8636f8732409 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260203192932-546029d2fa20 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace google.golang.org/genproto => google.golang.org/genproto/googleapis/rpc v0.0.0-20260203192932-546029d2fa20
