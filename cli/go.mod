module github.com/jwenz723/podlifecycle/cli

go 1.13

replace github.com/jwenz723/podlifecycle/proto => ../proto

require (
	github.com/jwenz723/podlifecycle/proto v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/sys v0.0.0-20190412213103-97732733099d // indirect
	google.golang.org/grpc v1.31.0
)
