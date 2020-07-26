APP = "gRPC"

run:
	go fmt ./... 

build:
	go build -o ${APP} .

server:
	go run main.go gs

client:
	go run main.go gc

proto:
	protoc -I data/protos/ data/protos/greet.proto --go_out=plugins=grpc:data/protos/greet