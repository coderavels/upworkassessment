setup:
	go mod tidy -v

start: build
	./virtualbookstore

build:
	go build -o virtualbookstore

generate:
	go generate ./... 

test:
	go test ./...
