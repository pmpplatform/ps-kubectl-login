run:
	go run main.go

all: prep binaries docker

prep:
	mkdir -p bin

binaries: linux64 darwin64

linux64:
	GOOS=linux GOARCH=amd64 go build -o bin/kubectl-login-linux .

darwin64:
	GOOS=darwin GOARCH=amd64 go build -o bin/kubectl-login-darwin .

pack: pack-linux64 pack-darwin64 

pack-linux64: linux64
	upx --brute bin/kubectl-login-linux

pack-darwin64: darwin64
	upx --brute bin/kubectl-login-darwin
