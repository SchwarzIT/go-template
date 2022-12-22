// +build tools

package main

import (
	// Go vulnerability scanner
	// https://go.dev/blog/vuln
	_ "golang.org/x/vuln/cmd/govulncheck"
	// golangci linter
	// https://golangci-lint.run
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	{{- if .Extensions.grpc.base }}
	// gRPC generation via https://buf.build/
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking"
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-lint"
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	{{- end }}
	{{- if .Extensions.grpc.grpcGateway }}
	// gRPC gateway
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/google/gnostic/cmd/protoc-gen-openapi"
	{{- end }}
)
