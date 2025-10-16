module mygomessenger.com/mygomessenger/test_client

go 1.25.2

require (
	github.com/jmoiron/sqlx v1.4.0
	mygomessenger.com/mygomessenger/gen/go/messages v0.0.0
	mygomessenger.com/mygomessenger/gen/go/users v0.0.0
	github.com/lib/pq v1.10.9
	google.golang.org/grpc v1.76.0
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace mygomessenger.com/mygomessenger/gen/go/users => ../gen/go/users

replace mygomessenger.com/mygomessenger/gen/go/messages => ../gen/go/messages

