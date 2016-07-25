ARCH ?= darwin
BINARY_NAME := chaperone
BUILD_DATE := $$(date +%Y-%m-%d-%H:%M)
ORG_PATH = github.com/mgoodness
REPO_PATH = $(ORG_PATH)/$(BINARY_NAME)
BUILD_DATE_VAR := $(REPO_PATH)/BuildDate
IMAGE_NAME := quay.io/mgoodness/$(BINARY_NAME)
REPO_VERSION := $$(git describe --abbrev=0 --tags)
GIT_HASH := $$(git rev-parse --short HEAD)
GIT_VAR := $(REPO_PATH)/GitCommit
VERSION_VAR := $(REPO_PATH)/Version
GOBUILD_VERSION_ARGS := --ldflags "-s -X $(VERSION_VAR)=$(REPO_VERSION) -X $(GIT_VAR)=$(GIT_HASH) -X $(BUILD_DATE_VAR)=$(BUILD_DATE)"

build: *.go fmt
	go build -o build/bin/$(ARCH)/$(BINARY_NAME) $(GOBUILD_VERSION_ARGS) \
		$(ORG_PATH)/$(BINARY_NAME)

clean:
	rm -rf build/bin/*
	-docker rm $(docker ps -a -f 'status=exited' -q)
	-docker rmi $(docker images -f 'dangling=true' -q)

cross:
	CGO_ENABLED=0 GOOS=linux go build -o build/bin/linux/$(BINARY_NAME) \
		$(GOBUILD_VERSION_ARGS) -a --ldflags '-extldflags "-static"' \
		$(ORG_PATH)/$(BINARY_NAME)

docker: cross
	docker build -t $(IMAGE_NAME):$(GIT_HASH) .

fmt:
	gofmt -w=true -s $$(find . -type f -name '*.go' -not -path "./vendor/*")
	goimports -w=true -d $$(find . -type f -name '*.go' -not -path "./vendor/*")

release: test docker
	docker push $(IMAGE_NAME):$(GIT_HASH)
	docker tag $(IMAGE_NAME):$(GIT_HASH) $(IMAGE_NAME):$(REPO_VERSION)
	docker push $(IMAGE_NAME):$(REPO_VERSION)

test:
	go test $$(glide nv)

version:
	@echo $(REPO_VERSION)

.PHONY: build
