// Package cmd is the command surface of odin cli tool provided by kubuskotak.
// # This manifest was generated by ymir. DO NOT EDIT.
//go:build tools
// +build tools

package cmd

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
