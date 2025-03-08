// Package proto ...
package proto

//go:generate protoc -I . --go-orb_out=paths=source_relative:. --go-orb_opt=supported_servers=grpc ./httpgateway_v1/httpgateway.proto
