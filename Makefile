GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

.PHONY: build
build:
	$(GOBUILD) -o gomodctl -v cmd/gomodctl/gomodctl.go

.PHONY: lint
lint:
	golint -set_exit_status=1 `go list ./...`

.PHONY: test
test:
	$(GOTEST) -v ./...
