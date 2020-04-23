# gRPC service + client boilerplate / scaffolding

Project status badges


[![CircleCI](https://circleci.com/gh/sunsingerus/mservice.svg?style=svg)](https://circleci.com/gh/sunsingerus/mservice)
[![issues](https://img.shields.io/github/issues/sunsingerus/mservice.svg)](https://github.com/sunsingerus/mservice/issues)
[![tags](https://img.shields.io/github/tag/sunsingerus/mservice.svg)](https://github.com/sunsingerus/mservice/tags)
[![Go Report Card](https://goreportcard.com/badge/github.com/sunsingerus/mservice)](https://goreportcard.com/report/github.com/sunsingerus/mservice)

## What is this
This is a gRPC client+service boilerplate / scaffolding. 
It exposes bi-directional stream service which consumes text file and upper-cases it.
 

## Most interesting parts are:
- `DataChunkFile` - ordered stream of data chunks with start/stop marks. 
  Used to transfer custom-sized data (possibly accompanied by metadata) over gRPC stream. 
  Inspired by `os.File`, implements `io.Writer`, `io.WriterTo`, `io.ReaderFrom`, `io.Closer` interfaces and thus is compatible/applicable in such functions as `io.Copy(dst, src)`.
  Located in [pkg/api/mservice/type_data_chunk_file.go](pkg/api/mservice/type_data_chunk_file.go). 
- Client server-less tests, based on `mockgen`-generated server-side mocks.
  Located in [pkg/controller/client/client_test.go](pkg/controller/client/client_test.go).
- Service tests
  - Client-less tests, used to test both `DataChunk` chunker/tansfer/aggregator and sever-side functionality
  - Network-less round-trip tests, used to test whole round-trip communication, with full-blown Server, launched during test case and Client dialing to Server.
    Based on custom in-memory dialer, network not used.
     
  Located in [pkg/controller/service/control_plane/server_test.go](pkg/controller/service/control_plane/server_test.go).
  Service tests allow to both test service and transport layer functionality.
- Complex `protobuf` nested messages with optional fields. accompanied by wrapper `Set*` functions.
  Wrapper functions are helpful, because generator generate `Get*` functions, but omit `Set*` functions, which is not convenient for messages with multiple optional fields. 
  Located in [pkg/api/mservice](pkg/api/mservice).
   

---
## Additional reading

- [How to run client and server parts][run]
- [How to install `protoc` compiler][protoc]

[protoc]: ./docs/protoc.md
[run]: ./docs/run.md
