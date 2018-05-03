BINARY := defang

VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || (git describe --always --long --dirty|tr '\n' '-';date +%Y.%m.%d))
LDFLAGS = -ldflags "-w -s -X main.version=${VERSION}"
LDFLAGS_DEV = -ldflags "-X main.version=${VERSION}"

#Build release builds
release: 
	@CGO_ENABLED=0 gox -osarch="darwin/386 darwin/amd64 linux/386 linux/amd64 windows/386 windows/amd64" ${LDFLAGS} -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" github.com/jakewarren/defang/cmd/defang

#Build a development build
dev: 
	@CGO_ENABLED=0 go build --installsuffix cgo ${LDFLAGS_DEV} -o bin/${BINARY} cmd/defang/main.go

#Install a release build on your local system
install: clean
	@go install ${LDFLAGS} github.com/jakewarren/defang/cmd/defang

build: 
	@go build ${LDFLAGS_DEV} -o bin/${BINARY} cmd/defang/main.go

clean: 
	@go clean -i

test:
	@go test -v ./...

# update the golden files used for the integration tests
update-tests:
	@go test integration/cli_test.go -update
