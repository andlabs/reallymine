PROJECT ?= $(notdir $(patsubst %/,%,$(CURDIR)))
PACKAGE ?= github.com/undeadbanegithub/$(PROJECT)

GO_SOURCES = $(shell find . -type f \( -iname '*.go' \) \
	-not \( -path "./vendor/*" -path ".*" \) \
	-not \( -path "./_notes/*" -path ".*" \) \
	-not \( -path "./_pre/*" -path ".*" \) \
	-not \( -path "./_sectors/*" -path ".*" \) \
	-not \( -path "./_try/*" -path ".*" \))

export GO111MODULE = on

fmt: $(GO_SOURCES)	# format go sources
	gofmt -w -s -l $^

.build/$(PROJECT).%: $(GO_SOURCES) go.mod
	mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=$(basename $*) GOARCH=$(patsubst .%,%,$(suffix $*)) go build -o $@ $(PACKAGE)

.BIN_LOCAL = $(PROJECT).$(shell go env GOOS).$(shell go env GOARCH)

build: .build/$(.BIN_LOCAL)
