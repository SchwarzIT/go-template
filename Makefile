SHELL=/bin/bash -e -o pipefail
PWD = $(shell pwd)

# constants
GOLANGCI_VERSION = 1.45.2

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
	@go run cmd/dotembed/main.go -target _template -o embed_gen.go -pkg gotemplate -var FS
	@go run cmd/options2md/main.go -o docs/options.md

GOLANGCI_LINT = bin/golangci-lint-$(GOLANGCI_VERSION)
$(GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b bin v$(GOLANGCI_VERSION)
	@mv bin/golangci-lint "$(@)"

lint: fmt $(GOLANGCI_LINT) download ## Lints all code with golangci-lint
	@$(GOLANGCI_LINT) run

lint-reports: out/lint.xml

.PHONY: out/lint.xml
out/lint.xml: $(GOLANGCI_LINT) out download
	$(GOLANGCI_LINT) run ./... --out-format checkstyle | tee "$(@)"

test: ## Runs all tests
	@go test ./...

coverage: out/report.json ## Displays coverage per func on cli
	go tool cover -func=out/cover.out

html-coverage: out/report.json ## Displays the coverage results in the browser
	go tool cover -html=out/cover.out

test-reports: out/report.json

.PHONY: out/report.json
out/report.json: out
	go test ./... -coverprofile=out/cover.out --json | tee "$(@)"

clean-test-project: ## Removes test-project
	@rm -rf testing-project

clean-test-project-cases: ## Removes test-project cases
	@rm -rf testing-project-*

clean: clean-test-project ## Cleans up everything
	@rm -rf bin out

ci: lint-reports test-reports

create-test-project: clean-test-project testing-project ## Creates a testing-project from the template

.PHONY: testing-project
testing-project:
	go run cmd/gt/*.go new -c pkg/gotemplate/testdata/values.yml

.PHONY: testing-projects-cases
testing-projects-cases: clean-test-project-cases ## test generating mutiple projects with running tests & linting within the new project
	@for case in basic grpc grpc-gateway ; do \
		echo Test: $$case ; \
		go run cmd/gt/*.go new -c test-cases/test-$$case.yml ; \
		cd testing-project-$$case ; \
		make ci ; \
		make docker ; \
		make all ; \
		cd .. ; \
    done

help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ''
