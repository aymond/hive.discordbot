PACKAGE ?= hive.discordbot
VERSION ?= latest
K8S_DIR       ?= ./k8s
K8S_BUILD_DIR ?= ./build_k8s
K8S_FILES     := $(shell find $(K8S_DIR) -name '*.yml' -or -name '*.yaml' | sed 's:$(K8S_DIR)/::g') 

DOCKER_REGISTRY_DOMAIN ?= docker.io
DOCKER_REGISTRY_PATH   ?= aymon
DOCKER_IMAGE           ?= $(DOCKER_REGISTRY_PATH)/$(PACKAGE):$(VERSION)
DOCKER_IMAGE_DOMAIN    ?= $(DOCKER_REGISTRY_DOMAIN)/$(DOCKER_IMAGE)

MAKE_ENV += PACKAGE VERSION DOCKER_IMAGE DOCKER_IMAGE_DOMAIN

SHELL_EXPORT := $(foreach v,$(MAKE_ENV),$(v)='$($(v))' )

default: build_in_docker ## build docker by default

build_in_docker:   ## build in docker
	rm -rfv bin
	docker build . -t "$(DOCKER_IMAGE)"

push-docker: build-docker
	docker push "$(DOCKER_IMAGE)"

fmt:  ## format all golang files
	go fmt

build:  ## build
	go build -o ./bin/main

$(K8S_BUILD_DIR):
	@mkdir -p $(K8S_BUILD_DIR)

help: ## prints out the menu of command options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build-k8s: $(K8S_BUILD_DIR)
	@for file in $(K8S_FILES); do \
		mkdir -p `dirname "$(K8S_BUILD_DIR)/$$file"` ; \
		$(SHELL_EXPORT) envsubst <$(K8S_DIR)/$$file >$(K8S_BUILD_DIR)/$$file ;\
	done

deploy: build-k8s push-docker # deploy
	kubectl apply -f $(K8S_BUILD_DIR)

.PHONY: default help build_in_docker build build-k8s deploy