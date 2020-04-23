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

Having these done correctly, we'll be able to compile with `protoc` files with `include` statements, like the following:
```.proto
import "google/protobuf/timestamp.proto";
```
