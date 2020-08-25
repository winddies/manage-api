run:
	go run ./main.go --config=./config.yml
build:
	GOOS=linux GOARCH=amd64 go build -o exec ./main.go
