all: mac

mac:
	env GOOS=darwin GOARCH=386 go build -o gethinx-darwin-386
linux:

