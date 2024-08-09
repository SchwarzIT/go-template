clean-test-project: ## Removes test-project
	@rm -rf testing-project

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
