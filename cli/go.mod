module github.com/jwenz723/podlifecycle/cli

go 1.13

replace github.com/jwenz723/podlifecycle/server => ./../server

require (
	github.com/jwenz723/podlifecycle/server v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0 // indirect
)
