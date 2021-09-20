# TODO: add goreleaser
# TODO: ginkgo support?

SHELL=/bin/bash -e -o pipefail
PWD = $(shell pwd)

# constants
GOLANGCI_VERSION = {{.golangciVersion}}
DOCKER_REPO = {{.appName}}
DOCKER_TAG = latest

bin/golangci-lint: bin/golangci-lint-$(GOLANGCI_VERSION)
	@ln -sf golangci-lint-$(GOLANGCI_VERSION) $@

bin/golangci-lint-$(GOLANGCI_VERSION):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b bin $(GOLANGCI_VERSION)
	@mv bin/golangci-lint $@

out:
	@mkdir -p out

out/bin:
	@mkdir -p out/bin

.PHONY: git-hooks
git-hooks: ## Binds the defined hooks from local .githooks directory to git config
	@git config --local core.hooksPath .githooks/

.PHONY: all
all: git-hooks tidy {{if .grpcEnabled}}generate{{end}} ## Initializes all tools
	@go mod tidy # hack to run tidy again after generating

.PHONY: env-%
env-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

.PHONY: download
download: ## Downloads the dependencies
	go mod download

.PHONY: tidy
tidy: ## Cleans up go.mod and go.sum
	go mod tidy

.PHONY: vendor
vendor: ## Vendors the packages to ./vendor
	go mod vendor

.PHONY: test-build
test-build: download ## Tests whether the code compiles
	go build -o /dev/null ./...

.PHONY: build
build: download out/bin ## Builds all binaries
	CGO_ENABLED=0 go build -ldflags="-w -s" -o out/bin ./...

.PHONY: build-linux
build-linux: download out/bin ## Builds all binaries for linux
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-w -s" -o out/bin ./...

.PHONY: fmt
fmt: ## Formats all code with go fmt
	@go fmt ./...

.PHONY: release
release:
	# TODO: goreleaser

.PHONY: lint
lint: bin/golangci-lint download ## Lints all code with golangci-lint
	bin/golangci-lint run

.PHONY: lint-reports
lint-reports: bin/golangci-lint out download ## Lints all code with golangci-lint and generates report
	bin/golangci-lint run ./... --out-format checkstyle > out/lint.xml

.PHONY: fix
fix: bin/golangci-lint ## Fix lint violations
	bin/golangci-lint run --fix

.PHONY: test
test: download ## Runs all tests
	go test ./...

.PHONY: test-reports
test-reports: out download ## Runs all tests and generates reports
	go test ./... -coverprofile=out/cover.out --json |tee out/report.json

.PHONY: clean-bin
clean-bin: ## Cleans local binary folders
	@rm -rf bin testbin

.PHONY: clean-outputs
clean-outputs: ## Cleans output folders out, vendor
	@rm -rf out vendor api/proto/google api/proto/validate {{if .grpcGatewayEnabled}}api/proto/protoc-gen-openapiv2{{end}}

.PHONY: clean-go
clean-go: ## Cleans module file
	@go mod tidy

.PHONY: clean
clean: clean-bin clean-outputs clean-go ## Cleans everything up

.PHONY: docker
docker: build-linux ## Builds docker image
	docker build -t $(DOCKER_REPO):$(DOCKER_TAG) .

{{if .grpcEnabled}}
# Go dependencies versioned through tools.go
GO_DEPENDENCIES = google.golang.org/protobuf/cmd/protoc-gen-go \
				google.golang.org/grpc/cmd/protoc-gen-go-grpc \
				github.com/envoyproxy/protoc-gen-validate \
				github.com/bufbuild/buf/cmd/buf \
                github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking \
                github.com/bufbuild/buf/cmd/protoc-gen-buf-lint
{{if .grpcGatewayEnabled}}
# additional dependencies for grpc-gateway
GO_DEPENDENCIES += github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
				github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
{{end}}

define make-go-dependency
  # target template for go tools, can be referenced e.g. via /bin/<tool>
  bin/$(notdir $1):
	GOBIN=$(PWD)/bin go install $1
endef

# this creates a target for each go dependency to be referenced in other targets
$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency, $(dep))))

api/proto/google:
	@mkdir -p api/proto/google/api
	curl -s https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto -o api/proto/google/api/annotations.proto
	curl -s https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto -o api/proto/google/api/http.proto
	curl -s https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto -o api/proto/google/api/field_behavior.proto

api/proto/validate:
	@mkdir -p api/proto/validate
	curl -s https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/main/validate/validate.proto -o api/proto/validate/validate.proto

{{if .grpcGatewayEnabled}}
api/proto/protoc-gen-openapiv2:
	@mkdir -p api/proto/protoc-gen-openapiv2/options
	curl -s https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/protoc-gen-openapiv2/options/openapiv2.proto -o api/proto/protoc-gen-openapiv2/options/openapiv2.proto
	curl -s https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/protoc-gen-openapiv2/options/annotations.proto -o api/proto/protoc-gen-openapiv2/options/annotations.proto
{{end}}

.PHONY: protolint
protolint: bin/buf bin/protoc-gen-buf-lint ## Lints your protobuf files
	bin/buf lint

.PHONY: protobreaking
protobreaking: bin/buf bin/protoc-gen-buf-breaking ## Compares your current protobuf with the version on master to find breaking changes
	bin/buf breaking --against '.git#branch=master'

.PHONY: generate
generate: ## Generates code from protobuf files
generate: api/proto/google api/proto/validate {{if .grpcGatewayEnabled}}api/proto/protoc-gen-openapiv2 bin/protoc-gen-grpc-gateway bin/protoc-gen-openapiv2{{end}} bin/buf bin/protoc-gen-go bin/protoc-gen-go-grpc bin/protoc-gen-validate
	PATH=$(PWD)/bin:$$PATH buf generate --path api/proto/{{.appName}}/v1/{{.appName}}.proto
{{end}}

.PHONY: ci
ci: lint-reports test-reports ## Executes lint and test and generates reports

.PHONY: help
help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''
