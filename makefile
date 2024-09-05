gotest:
	go test ./...

runServer:
	go run main.go

build:
	go build -o data_recovery main.go 

build-cli:
	go build -o client client/client.go

cli:
	go run client/client.go -d http://localhost:8888/test/jdk-8u361-windows-x64.exe

