GOIMPORTS_BIN ?= $(shell which goimports || true)
GOTEST_FLAGS ?= -v --coverprofile main.coverprofile --covermode count --cover github.com/aperturerobootics/cli/cmd/urfave-cli-genflags
GOBUILD_FLAGS ?= -x

export GOIMPORTS_BIN

.PHONY: all
all: test build smoke-test

.PHONY: test
test:
	go test $(GOTEST_FLAGS) ./...

.PHONY: build
build:
	go build $(GOBUILD_FLAGS) ./...

.PHONY: smoke-test
smoke-test: build
	./urfave-cli-genflags --help

.PHONY: show-cover
show-cover:
	go tool cover -func main.coverprofile

.PHONY: run
run: build
	./urfave-cli-genflags	
	# ./urfave-cli-genflags -f altsrc/flag-spec.yaml -o altsrc/flag_generated.go -p altsrc -a
