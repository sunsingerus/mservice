# gRPC service + client

Project status


[![CircleCI](https://circleci.com/gh/sunsingerus/mservice.svg?style=svg)](https://circleci.com/gh/sunsingerus/mservice)
[![issues](https://img.shields.io/github/issues/sunsingerus/mservice.svg)](https://github.com/sunsingerus/mservice/issues)
[![tags](https://img.shields.io/github/tag/sunsingerus/mservice.svg)](https://github.com/sunsingerus/mservice/tags)
[![Go Report Card](https://goreportcard.com/badge/github.com/sunsingerus/mservice)](https://goreportcard.com/report/github.com/sunsingerus/mservice)

## How to install `protoc`

- Download latest protobuf release from [here](https://github.com/protocolbuffers/protobuf/releases)
- We'll have something like `protoc-3.11.4-linux-x86_64.zip` with the following structure:
```text
    bin
        protoc
    include
        google
            protobuf
                ... many files here ...
```
- Place `bin` into `$PATH`-searchable - `bin`
- Place `include` near `bin`, so we'll have something like the following:
```text
    bin
        ... $PATH-searchable bin folder ...
        ... you may have your old bin files ...
        protoc
    include
        google
            protobuf
                ... many files here ...
``` 

Having these done correctly, we'll be avle to compile with `protoc` files with `include` statements, like the following:
```.proto
import "google/protobuf/timestamp.proto";
```
