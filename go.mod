module github.com/ItakawaM/gophpipe

go 1.27.3

require go.nanomsg.org/mangos/v3 v3.6.0 // direct

require (
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.6.2 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

tool (
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)
