run:
	go run main.go

build:
	GOOS=windows GOARCH=amd64 go build -o gproxy.exe main.go
	go build -o gproxy main.go

