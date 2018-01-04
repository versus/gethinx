all: mac

mac:
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build  -a --installsuffix cgo --ldflags="-s" -o gethinx-darwin-x64

linux:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o gethinx-linux-x64

clean:
	rm gethinx-darwin-x64 || true
	rm gethinx-linux-x64 || true
