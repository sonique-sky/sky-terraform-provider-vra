CIRCLE_BUILD_NUM ?= 0

APP=sky-terraform-provider-vra
PKG=/go/src/github.com/sonique-sky/${APP}
TAG=timwoolford/${APP}:0.1.$(CIRCLE_BUILD_NUM)

BIN=$(firstword $(subst :, ,${GOPATH}))/bin
VENDOR = $(BIN)/govendor
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: gobuild
gobuild: vendor ; $(info $(M) building…)
	GOOS=darwin go build -v -o bin/${APP}

.PHONY: gotest
gotest: gobuild ; $(info $(M) running tests…)
	@go test ./...

.PHONY: build
build:
	docker run --rm \
	 -v "${PWD}":${PKG} \
	 -w ${PKG} \
	 golang:1.10	 \
	 make gobuild

.PHONY: build-image
build-image:
	docker build -t ${TAG} .

.PHONY: push-image
push-image:
	docker push ${TAG}

.PHONY: clean
clean: ; $(info $(M) cleaning…)
	@docker images -q ${APP} | xargs docker rmi -f
	@rm -rf bin/*

vendor:
	command -v $(VENDOR) >/dev/null 2>&1 || go get -u github.com/kardianos/govendor
	$(VENDOR) init
	$(VENDOR) add +external
	@touch $@
