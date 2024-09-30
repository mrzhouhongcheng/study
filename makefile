run:
	go run main.go

build:
	go build -o gproxy main.go
	GOOS=windows GOARCH=amd64 go build -o gproxy.exe main.go

