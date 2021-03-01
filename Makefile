run:
	go run main.go

all: prep binaries docker

prep:
	mkdir -p bin

binaries: linux-amd64 darwin-amd64 darwin-arm64

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/kubectl-login-linux-amd64 .

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/kubectl-login-darwin-amd64 .

darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o bin/kubectl-login-darwin-arm64 .

pack: pack-linux-amd64 pack-darwin-amd64 #pack-darwin-arm64 

pack-linux-amd64: linux-amd64
	upx --brute bin/kubectl-login-linux-amd64

pack-darwin-amd64: darwin-amd64
	upx --brute bin/kubectl-login-darwin-amd64

#pack-darwin-arm64: darwin-arm64
#	upx --brute bin/kubectl-login-darwin-arm64
