statik:
	statik -m -src=static/content -dest=static

go-build:
	go build

LDFLAGS=-ldflags="-buildid= -X 'github.com/jbrunton/g3ops/cmd.Version=${version}'"
go-build-release:
	@test -n "$(version)" || (echo '$$version required' && exit 1)
	export CGO_ENABLED=0
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -trimpath -o g3ops-darwin-amd64
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -trimpath -o g3ops-linux-amd64

compile: statik go-build

compile-release: statik go-build-release

build: compile test

unit-test:
	go test -coverprofile c.out $$(go list ./... | grep -v /e2e)

test: unit-test

.PHONY: statik go-build go-build-release compile compile-release build unit-test test

.DEFAULT_GOAL := build
