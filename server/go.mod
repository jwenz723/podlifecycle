module github.com/jwenz723/podlifecycle/server

go 1.13

replace github.com/jwenz723/podlifecycle/proto => ../proto

require (
	github.com/jwenz723/podlifecycle/proto v0.0.0-00010101000000-000000000000
	github.com/oklog/run v1.1.0
	go.uber.org/zap v1.15.0
	google.golang.org/grpc v1.31.0
)
