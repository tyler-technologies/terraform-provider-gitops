GOPATH := $(shell go env | grep GOPATH | sed 's/GOPATH="\(.*\)"/\1/')
PATH := $(GOPATH)/bin:$(PATH)
export $(PATH)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
TEST_DESTS := $(dir $(wildcard ./test/*/*_test.go))

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

fetch: ## download makefile dependencies
	@hash goreleaser 2>/dev/null || go get -u -v github.com/goreleaser/goreleaser

clean: ## cleans previously built binaries
	rm -rf ./dist
	@cd test && $(MAKE) clean

publish: clean fetch ## publishes assets
	@if [ "${GITHUB_TOKEN}" == "" ]; then\
	  echo "GITHUB_TOKEN is not set";\
		exit 1;\
	fi
	@if [ "$(GIT_BRANCH)" != "master" ]; then\
	  echo "Current branch is: '$(GIT_BRANCH)'.  Please publish from 'master'";\
		exit 1;\
	fi
	git tag -a $(VERSION) -m "$(MESSAGE)"
	git push --follow-tags
	$(GOPATH)/bin/goreleaser

build: clean fetch ## publishes in dry run mode
	$(GOPATH)/bin/goreleaser release --snapshot --skip-validate --skip-publish --skip-sign


.PHONY: test copyplugins

copyplugins: ## copy plugins to test folders
	$(eval OS_DIRS := $(dir $(wildcard ./dist/terraform-provider-gitops*/*)))
	$(eval OS_ARCH := $(patsubst ./dist/terraform-provider-gitops_%/, %, $(OS_DIRS)))
	@sleep 1
	@for f in $(TEST_DESTS); do \
		for o in $(OS_ARCH); do \
	  	mkdir -p $$f/terraform.d/plugins/$$o; \
			cp ./dist/terraform-provider-gitops_$$o/* $$f/terraform.d/plugins/$$o; \
		done; \
	done

test: copyplugins ## test
	@cd test && $(MAKE) test

fulltest: build test ## build and test