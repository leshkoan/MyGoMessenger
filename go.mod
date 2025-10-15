module github.com/leshkoan/MyGoMessenger

go 1.25.2

require (
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	github.com/segmentio/kafka-go v0.4.49
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)

replace github.com/leshkoan/MyGoMessenger/gen/go/users => ./gen/go/users

replace github.com/leshkoan/MyGoMessenger/gen/go/messages => ./gen/go/messages
