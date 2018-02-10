name = versus197/gethinx
version ?= latest

all: mac

# This makefile contains some convenience commands for deploying and publishing.

# For example, to build docker container locally, just run:
# $ make docker

# or to docker  images  the :latest version to the specified registry as :1.0.0, run:
# $ make docker version=1.0.0


dist: clean
	mkdir artefacts
	cp -r ./templates artefacts
	cp config.toml artefacts/config.toml

build: dist
	docker run --rm -v "${GOPATH}":/gopath -v "$(CURDIR)":/app -e "GOPATH=/gopath" -e "GOPATH=/gopath" -w /app golang:1.9 sh -c 'go build -a --installsuffix cgo --ldflags="-s"    -o ./artefacts/gethinx-linux-x64'


mac: dist
	$(call blue, "Building MacOS binary...")
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build  -a --installsuffix cgo --ldflags="-s" -o ./artefacts/gethinx-darwin-x64

linux: dist
	$(call blue, "Building Linux binary...")
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -a --installsuffix cgo --ldflags="-s" -o ./artefacts/gethinx-linux-x64

docker: dist
	$(call blue, "Building docker image...")
	docker build -t ${name}:${version} .

clean:
	rm -rf artefacts


define blue
	@tput setaf 6
	@echo $1
	@tput sgr0
endef
