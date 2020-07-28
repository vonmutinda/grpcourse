## gRPC - MasterClass
- A code along repo for [Udemy's gRPC [Golang] Masterclass course](https://www.udemy.com/course/grpc-golang/)

## Description
`gRPC` is a `RPC` framework developed by Google for building high performance APIs. It's uses [Protocol Buffers](https://developers.google.com/protocol-buffers) as the interface definition language. gRPC payloads are significantly smaller than a JSON equivalents hence faster :- low latency over the wire. A good alternative for REST API. 

## Concepts Explored
- [x] Unary Streaming 
- [x] Server Streaming
- [x] Client Streaming  
- [x] Bi-directional Streaming 
- [x] Error Handling & RPC deadlines
- [x] Security - SSL Authentication
- [x] File Upload - (client streaming)
- [ ] CRUD with MongoDB 

### Local Set Up
This code is entirely educational and definitely not production-ready 
+ Clone the app `git clone https://github.com/vonmutinda/grpcourse.git` 
+ Install `MongoDB` databasase. 
+ Inspect the `Makefile` for the various commands of interracting with the app. 
+ Run `make server` and `make client` on two seperate tabs.;


## Technologies Used 
A list of technologies used in this project:
- [Golang version `go1.14.6`](https://golang.org) 
- [gRPC](https://github.com/grpc/grpc-go) A RPC framework that leverages on `HTTP/2` protocol which has powerful multiplexing capabilites.
- [MongoDB](https://www.mongodb.com/try/download/community) Community Server

## Resources (Further Read)
Some of the resources that I found useful `Go/Golang`. 
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [gRPC](https://grpc.io/)
- [Building Microservices in Go](https://www.youtube.com/playlist?list=PLmD8u-IFdreyh6EUfevBcbiuCKzFk0EW_) - Nic Jackson [Start Here](https://www.youtube.com/watch?v=pMgty_RYIOc&list=PLmD8u-IFdreyh6EUfevBcbiuCKzFk0EW_&index=14&t=0s)
- [The Complete gRPC Course [Protobuf, Go, Java]](https://www.youtube.com/playlist?list=PLy_6D98if3UJd5hxWNfAqKMr15HZqFnqf) : An Advanced gRPC Course
