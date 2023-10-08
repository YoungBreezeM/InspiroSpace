.PHONY: all build docker-build

all:deploy

build:gen 
	GOOS=linux GOARCH=amd64 go build -o build/freebox -v -ldflags="-s -w" cmd/main.go && upx -9 build/freebox

gen:
	go generate ./...

check:
	go build -o build/main -gcflags=-m cmd/main.go

docker-build:
	scp build/freebox prod:~/ops-pkg/apps/freebox
	ssh prod "cd ops-pkg/apps; ./build-image freebox"

deploy:build docker-build
	ssh prod "cd ops-pkg/apps; ./deploy freebox"

clean:
	rm -f build/*