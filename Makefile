APPNAME = azureip2tfvars
VERSION ?= latest
LDFLAGS = -ldflags "-X main.VERSION=$(VERSION) -X main.COMMIT=$(shell git rev-parse --short HEAD) -X main.BRANCH=$(shell git branch | grep \* | cut -d ' ' -f2)"

default: clean macos-build windows-build linux-build

clean:
	rm -fr bin/

macos-build: dep
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/$(APPNAME)-darwin-amd64

windows-build: dep
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/$(APPNAME)-windows-amd64

linux-build: dep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/$(APPNAME)-linux-amd64

dep:
	go get ./...

test:
	go test -v ./...

coverage:
	go test -cpu=1 -v ./... -failfast -coverprofile=coverage.txt -covermode=count
