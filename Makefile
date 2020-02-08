GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

.PHONY: build
build:
	$(GOBUILD) -o gomodctl -v cmd/gomodctl/gomodctl.go

.PHONY: lint
lint:
	golint -set_exit_status=1 `go list ./...` | grep -v tools

.PHONY: test
test:
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic -tags=integration ./...

.PHONY: test
download:
	@echo Download go.mod dependencies
	@go mod download

.PHONY: test
install-tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %