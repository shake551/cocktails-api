GOOS ?= linux
GOARCH ?= amd64

ROOT_PACKAGE := github.com/shake551/cocktails-api
SRC := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
OPEN_CHAT_APP_SERVER_BIN := ./bin/cocktails-api-server
CREATE_DUMMY_DATA_BIN := ./create_dummy_data

$(OPEN_CHAT_APP_SERVER_BIN): $(SRC)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -tags netgo -installsuffix netgo -ldflags '-s -w -extldflags "-static"' -o $@

$(CREATE_DUMMY_DATA_BIN): $(SRC)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -tags netgo -installsuffix netgo -ldflags '-s -w -extldflags "-static"' -o $@ $(ROOT_PACKAGE)/cmd/create_dummy_data

.PHONY: build
build: $(OPEN_CHAT_APP_SERVER_BIN)

.PHONY: build_create_dummy_data
build_create_dummy_data: $(CREATE_DUMMY_DATA_BIN)

.PHONY: clean
clean:
	@rm -f $(OPEN_CHAT_APP_SERVER_BIN) $(CREATE_DUMMY_DATA_BIN)

.PHONY: check
check: fmt vet

.PHONY: vet
vet:
	@go vet $$(go list ./...)

.PHONY: fmt
fmt:
	@gofmt -l .