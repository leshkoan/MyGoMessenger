module test_client

go 1.25.2

require (
	github.com/leshkoan/MyGoMessenger/gen/go/users v0.0.0
	github.com/leshkoan/MyGoMessenger/gen/go/messages v0.0.0
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

replace github.com/leshkoan/MyGoMessenger/gen/go/users => ../gen/go/users

replace github.com/leshkoan/MyGoMessenger/gen/go/messages => ../gen/go/messages