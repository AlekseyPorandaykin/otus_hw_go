package main

//nolint
//go:generate protoc --go_out=./internal/server/grpc/pb --go_opt=paths=source_relative --go-grpc_out=./internal/server/grpc/pb --go-grpc_opt=paths=source_relative  api/grpc/EventService.proto
