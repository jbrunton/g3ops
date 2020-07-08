statik:
	statik -m -src=static-content

compile:
	go build

compile-release:
	export CGO_ENABLED=0
	repo_flags="-ldflags=-buildid= -trimpath"
	GOOS=darwin GOARCH=amd64 go build $$repo_flags -o g3ops-darwin-amd64
	GOOS=linux GOARCH=amd64 go build $$repo_flags -o g3ops-linux-amd64
	GOOS=windows GOARCH=amd64 go build $$repo_flags -o g3ops-windows-amd64.exe

build: statik compile

build-release: statik compile-release

.PHONY: statik compile build compile-release build-release
