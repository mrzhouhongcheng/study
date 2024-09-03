gotest:
	go test ./...

run:
	go run main.go

build:
	go build -o data_recovery main.go 

cli:
	go run client/client.go -d http://localhost:8888/test/jdk-8u361-windows-x64.exe

