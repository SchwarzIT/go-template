# TODO: add goreleaser
# TODO: ginkgo support?

SHELL=/bin/bash -e -o pipefail
PWD = $(shell pwd)

# constants
GOLANGCI_VERSION = 1.42.1
DOCKER_REPO = gt
DOCKER_TAG = latest

bin/golangci-lint: bin/golangci-lint-$(GOLANGCI_VERSION)
	@ln -sf golangci-lint-$(GOLANGCI_VERSION) $@

bin/golangci-lint-$(GOLANGCI_VERSION):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b bin v$(GOLANGCI_VERSION)
	@mv bin/golangci-lint $@

out:
	@mkdir -p out

out/bin:
	@mkdir -p out/bin

.git:
	@git init

go.mod:
	go mod init github.com/schwarzit/go-template

.PHONY: git-hooks
git-hooks: ## Binds the defined hooks from local .githooks directory to git config
	@git config --local core.hooksPath .githooks/

.PHONY: all
all: .git git-hooks bin/golangci-lint go.mod tidy download  ## Initializes all tools
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
	@rm -rf out vendor api/proto/google api/proto/validate 

.PHONY: clean-go
clean-go: ## Cleans module file
	@go mod tidy

.PHONY: clean
clean: clean-bin clean-outputs clean-go ## Cleans everything up

.PHONY: docker
docker: build-linux ## Builds docker image
	docker build -t $(DOCKER_REPO):$(DOCKER_TAG) .



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
