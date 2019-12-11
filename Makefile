.DEFAULT_GOAL := help

.PHONY: all
all:

.PHONY: build
build: ## Build for local environment
	@go build

.PHONY: run
run: build ## Run example script
	# stein loads the HCL files located on .policy directory by default
	# in addition, .policy directory can be overridden by each directory of given arguments
	#
	# in this case,
	#   stein applies rules located in these default directory to _examples/manifests/microservices/*/*/*/*
	#   * _examples/.policy/
	#   * _examples/manifests/.policy/
	#   stein doesn't apply this rules to them
	#   * _examples/spinnaker/.policy/
	#
	# Regardless of the default directory placed under the given path,
	# the following environment variables can be specified for the policy applied to all paths.
	# this variables can take multiple values separated by a comma, also can take directories and files
	#
	# export STEIN_POLICY=root-policy/,another-policy/special.hcl
	./stein apply \
		_examples/manifests/microservices/*/*/*/* \
		_examples/spinnaker/*/*/*

.PHONY: release
release: ## Build for multiple OSs, packaging it and upload to GitHub Release
	@bash ./scripts/release.sh

.PHONY: help
help: ## Show help message for Makefile target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: docs-build
docs-build: ## Build documentations with mkdocs
	@docker build -t mkdocs docs

.PHONY: docs-live
docs-live: build-docs ## Live viewing with mkdocs
	@docker run --rm -it -p 3000:3000 -v ${PWD}:/docs mkdocs

.PHONY: docs-deploy
docs-deploy: docs-build ## Deploy generated documentations to gh-pages
	@docker run --rm -it -v ${PWD}:/docs -v ~/.ssh:/root/.ssh mkdocs mkdocs gh-deploy

.PHONY: test
test: ## Run test
	@go test -v -race ./...
