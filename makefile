.DEFAULT_GOAL := all

GO=$(shell which go)
DISTVER=$(shell git describe --always --dirty --long --tags)
PKG=$(shell head -1 go.mod | sed 's/^module //')

all: test dist

test:
	$(GO) test -v ./scid/... ./csv/... ./util/...
#	$(GO) test -v ./...

dist:
	$(GO) build -ldflags "-X $(PKG)/util.Version=$(DISTVER)"

install: test dist
	$(GO) install -ldflags "-X $(PKG)/util.Version=$(DISTVER)" .

upgrade:
	$(GO) get -u && $(GO) mod tidy

goformat:
	gofmt -s -w .

