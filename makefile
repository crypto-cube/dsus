GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

all: build-debug

build-debug:
	$(GOBUILD) -v
clean:
	$(GOCLEAN)
	rm ./scripts/*.deb
build-deb: build-release
	cd scripts && /bin/bash genDEB.sh;
build-release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -trimpath -ldflags="-s -w -X main.isDebug=false" -a -v 
