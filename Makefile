SHELL=/bin/bash -e -o pipefail
PWD = $(shell pwd)
LINTER_VERSION = v1.55.1

all: git-hooks generate ## Initializes all tools and files

out:
	@mkdir -pv "$(@)"

test-build: ## Tests whether the code compiles
	@go build -o /dev/null ./...

build: out/bin ## Builds all binaries

GO_BUILD = mkdir -pv "$(@)" && go build -ldflags="-w -s" -o "$(@)" ./...
.PHONY: out/bin
out/bin:
	$(GO_BUILD)

git-hooks:
	@git config --local core.hooksPath .githooks/

download: ## Downloads the dependencies
	@go mod download

fmt: ## Formats all code with go fmt
	@go fmt ./...

run: fmt ## Run a controller from your host
	@go run ./main.go

generate: ## Generates files
	@go run cmd/options2md/main.go -o docs/options.md
	@go run github.com/nix-community/gomod2nix@latest --outdir nix


lint: fmt download ## Lints all code with golangci-lint
	@go run -v github.com/golangci/golangci-lint/cmd/golangci-lint@$(LINTER_VERSION) run

govulncheck: ## Vulnerability detection using govulncheck
	@go run golang.org/x/vuln/cmd/govulncheck ./...

test: ## Runs all tests
	@go test ./...

coverage: out/report.json ## Displays coverage per func on cli
	go tool cover -func=out/cover.out

html-coverage: out/report.json ## Displays the coverage results in the browser
	go tool cover -html=out/cover.out

test-coverage: out ## Creates a test coverage profile
	go test -v -cover ./... -coverprofile out/coverage.out -coverpkg ./...
	go tool cover -func out/coverage.out -o out/coverage.out

clean-test-project: ## Removes test-project
	@rm -rf testing-project

clean: clean-test-project ## Cleans up everything
	@rm -rf bin out

.PHONY: testing-project
testing-project: clean-test-project ## Creates a testing-project from the template
	@go run cmd/gt/*.go new -c $$VALUES_FILE

.PHONY: testing-project-ci-single
testing-project-ci-single:  ## Creates a testing-project from the template and run make ci within it
	@make testing-project VALUES_FILE=$$VALUES_FILE
	@make -C testing-project ci
	@make -C testing-project all

.PHONY: testing-project-default
testing-project-default: ## Creates the default testing-project from the template
	@make testing-project VALUES_FILE=pkg/gotemplate/testdata/values.yml

.PHONY: testing-project-ci
testing-project-ci:  ## Creates for all yml files in ./test_project_values a test project and run `make ci`
	for VALUES in ./test_project_values/*.yml; do \
		make testing-project-ci-single VALUES_FILE=$$VALUES; \
	done

.PHONY: release
release:  ## Create a new release version
	@./hack/release.sh

help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ''
