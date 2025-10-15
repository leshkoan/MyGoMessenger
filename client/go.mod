module client

go 1.25.2

require (
	github.com/leshkoan/MyGoMessenger/gen/go/messages v0.0.0
	github.com/leshkoan/MyGoMessenger/gen/go/users v0.0.0
	google.golang.org/grpc v1.76.0
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace github.com/leshkoan/MyGoMessenger/gen/go/users => ../gen/go/users

replace github.com/leshkoan/MyGoMessenger/gen/go/messages => ../gen/go/messages
