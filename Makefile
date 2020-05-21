ifeq ($(origin GITROOT), undefined)
GITROOT := $(shell git rev-parse --show-toplevel)
endif

# Commands
GO           ?= go
GOIMPORTS    ?= $(or $(wildcard $(GITROOT)/bin/goimports), $(shell which goimports))
GOLANCI_LINT ?= $(or $(wildcard $(GITROOT)/bin/golangci-lint), $(shell which golangci-lint))
GOLINT       ?= $(GOLANCI_LINT) run
PACKR        ?= $(or $(wildcard $(GITROOT)/bin/packr), $(shell which packr))
DOCKER       ?= docker
GREP         ?= grep
MKDIR_P      ?= mkdir -p

GIT_LATEST_TAG ?= $(shell git describe --tags --always HEAD)

# GO build flags
GO_IMPORT_PATH = github.com/haostudio/golinks
GO_BUILD_LDFLAGS =
GO_BUILD_LDFLAGS += -X $(GO_IMPORT_PATH)/internal/version.version=$(GIT_LATEST_TAG)

GO_BUILD_FULL_OPTS = -i $(GO_BUILD_OPTS) -ldflags "$(GO_BUILD_LDFLAGS)"

BUILDDIR ?= $(GITROOT)/build

# Docker parameters
DOCKER_REGISTRY       ?= haostudio
DOCKER_GOLINKS_CONFIG ?= local
DOCKER_TAG            ?= $(GIT_LATEST_TAG)-$(DOCKER_GOLINKS_CONFIG)
DOCKER_IMAGE          ?= $(DOCKER_REGISTRY)/golinks:$(DOCKER_TAG)

define INSTALL_RULE
install-$1:
ifeq (,$(shell which $1))
	$2
else
	@echo "$1 is installed"
endif
endef

define GO_INSTALL_RULE
$(eval $(call INSTALL_RULE,$1,GOBIN=$(GITROOT)/bin $(GO) get $2))
endef

define BUILD_RULE
$1: pre-build
	$(GO) build $(GO_BUILD_FULL_OPTS) -o $(BUILDDIR)/$1 $(GITROOT)/cmd/$1
endef

GO_DEPS = \
	packr;github.com/gobuffalo/packr/packr \
	goimports;golang.org/x/tools/cmd/goimports

APPS = \
	golinks

.PHONY: default pre-build wiki

default: wiki apps

deps: install-golangci-lint $(foreach DEP,$(GO_DEPS), $(eval CMD = $(word 1,$(subst ;, ,$(DEP)))) install-$(CMD))
$(foreach DEP,$(GO_DEPS), \
	$(eval CMD = $(word 1,$(subst ;, ,$(DEP)))) \
	$(eval SRC = $(word 2,$(subst ;, ,$(DEP)))) \
	$(eval $(call GO_INSTALL_RULE,$(CMD),$(SRC))) \
)

$(eval \
	$(call INSTALL_RULE,golangci-lint, \
	@wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.23.8)\
)

precommit: tidy apps format lint test

format:
	find . -name "*.go" | xargs $(GOIMPORTS) -w -local $(GO_IMPORT_PATH)

lint:
	$(GOLINT) $(GITROOT)/...

test:
	$(GO) test $(GITROOT)/... -short

full-test:
	$(GO) test $(GITROOT)/...

tidy:
	$(GO) mod tidy

pre-build:
	@$(MKDIR_P) $(BUILDDIR)
	$(PACKR)

clean:
	rm -rf $(BUILDDIR)/*
	$(PACKR) clean

apps: $(APPS)
$(foreach app, $(APPS), $(eval $(call BUILD_RULE,$(app))))

wiki:
	$(DOCKER) run --rm \
		-v $(GITROOT)/wiki:/docs \
		-v $(GITROOT)/images:/docs/docs/img \
		squidfunk/mkdocs-material build --clean

docker: wiki
	$(DOCKER) build --build-arg GOLINKS_CONFIG=$(DOCKER_GOLINKS_CONFIG) -t $(DOCKER_IMAGE) .

docker-push:
	$(DOCKER) push $(DOCKER_IMAGE)

jaeger:
	$(DOCKER) run -d --name jaeger \
	-e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
	-p 5775:5775/udp \
	-p 6831:6831/udp \
	-p 6832:6832/udp \
	-p 5778:5778 \
	-p 16686:16686 \
	-p 14268:14268 \
	-p 9411:9411 \
	jaegertracing/all-in-one:latest
