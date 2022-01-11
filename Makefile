all: build

build: mgr_service

mgr_service:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./mgr_service ./cmd/main.go
