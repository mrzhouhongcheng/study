gotest:
	go test ./...

runServer:
	go run main.go

build:
	go build -o data_recovery main.go 
	GOOS=windows GOARCH=amd64 go build -o data_recovery.exe main.go

build-cli:
	go build -o client client/client.go
	GOOS=windows GOARCH=amd64 go build -o client client/client.go

cli:
	go run client/client.go -d http://localhost:8888/test/jdk-8u361-windows-x64.exe

