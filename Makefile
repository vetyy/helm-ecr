PKG := github.com/vetyy/helm-ecr
GO111MODULE := on

.PHONY: all
all: deps build

.PHONY: deps
deps:
	@go mod download
	@go mod vendor
	@go mod tidy

.PHONY: build
build:
	@./hack/build.sh $(CURDIR) $(PKG)

.PHONY: install
install:
	@./hack/install.sh

.PHONY: lint
lint:
	@golangci-lint run ./...
