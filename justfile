gocmd := "go"

default: build-debug

build-debug:
    {{gocmd}} build -v

version := `git describe --tags --always 2>/dev/null || echo "dev"`

clean:
    {{gocmd}} clean

build-release:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 {{gocmd}} build -trimpath -ldflags="-s -w -X main.isDebug=false -X main.version={{version}}" -a -v
