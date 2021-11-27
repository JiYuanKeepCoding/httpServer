PACKAGE ?= hello-world
VERSION ?= 0.03

K8S_DIR ?= ./k8s
K8S_BUILD_DIR ?= ./build_k8s
K8S_FILES     := $(shell find $(K8S_DIR) -name '*.yml' | sed 's:$(K8S_DIR)/::g')

DOCKER_REGISTRY_DOMAIN ?= docker.io
DOCKER_REGISTRY_PATH   ?= jyabcd360/http_server
DOCKER_IMAGE           ?= $(DOCKER_REGISTRY_PATH)/$(PACKAGE):$(VERSION)
DOCKER_IMAGE_DOMAIN    ?= $(DOCKER_REGISTRY_DOMAIN)/$(DOCKER_IMAGE)
NAMESPACE              ?= http-server

MAKE_ENV += PACKAGE VERSION DOCKER_IMAGE DOCKER_IMAGE_DOMAIN NAMESPACE

SHELL_EXPORT := $(foreach v,$(MAKE_ENV),$(v)='$($(v))' )

.PHONY: build-docker
build-docker:
	docker build ./ -t "$(DOCKER_IMAGE)"

.PHONY: push-docker
push-docker: build-docker
	docker push "$(DOCKER_IMAGE)"

$(K8S_BUILD_DIR):
	@mkdir -p $(K8S_BUILD_DIR)

.PHONY: build-k8s
build-k8s: $(K8S_BUILD_DIR)
	@for file in $(K8S_FILES); do \
			mkdir -p `dirname "$(K8S_BUILD_DIR)/$$file"` ; \
			$(SHELL_EXPORT) envsubst <$(K8S_DIR)/$$file >$(K8S_BUILD_DIR)/$$file ;\
	done

.PHONY: deploy
deploy: build-k8s push-docker
	kubectl apply -f $(K8S_BUILD_DIR) -n ${NAMESPACE}