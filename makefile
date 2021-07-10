.DEFAULT_GOAL := dist

MAIN=sc-data-util

GO=$(shell which go)
DISTVER=$(shell git describe --always --dirty --long --tags)
PKG=github.com/RileyR387/sc-data-util

dist:
	$(GO) build -ldflags "-X main.Version=$(DISTVER) -X $(PKG)/scid.Version=$(DISTVER) -X $(PKG)/csv.Version=$(DISTVER) -X $(PKG)/util.Version=$(DISTVER)"

install:
	$(GO) install -ldflags "-X main.Version=$(DISTVER) -X $(PKG)/scid.Version=$(DISTVER) -X $(PKG)/csv.Version=$(DISTVER) -X $(PKG)/util.Version=$(DISTVER)" .

dev:
	$(GO) run -ldflags "-X main.Version=$(DISTVER) -X $(PKG)/scid.Version=$(DISTVER) -X $(PKG)/csv.Version=$(DISTVER) -X $(PKG)/util.Version=$(DISTVER)" . 2>>sc-dtc-client.err

test:
	$(GO) test -v ./...

race:
	$(GO) run -ldflags "-X main.Version=$(DISTVER) -X $(PKG)/scid.Version=$(DISTVER) -X $(PKG)/csv.Version=$(DISTVER) -X $(PKG)/util.Version=$(DISTVER)" --race . 2>>sc-dtc-client.err

goformat:
	gofmt -s -w .

