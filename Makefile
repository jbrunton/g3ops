statik:
	statik -src=static-content

compile:
	go build

build: statik compile
