# === Builder ===

FROM golang:1.12.7 AS builder

# Install required packages
RUN apt-get update && apt-get install -y apt-utils gettext-base

# Reconstruct source tree inside docker
WORKDIR /code
ADD . .

# ./vendor is excluded in .dockerignore, reconstruct it with 'mod' tool
RUN GO111MODULE=on go mod tidy
RUN GO111MODULE=on go mod vendor

# Before start build process, we need to ensure go test is OK
RUN ./dev/run_test.sh

# Build operator binary with explicitly specified output
RUN OUTPUT_BIN=/tmp/service ./dev/go_build_service.sh

# Build grpc-health-probe
RUN CGO_ENABLED=0 GO111MODULE=off go get -a github.com/grpc-ecosystem/grpc-health-probe

# === Image ===

FROM alpine:3.10 AS image

RUN apk add --no-cache ca-certificates

# For more convenient usage add bash
RUN apk add --no-cache bash

WORKDIR /

# Copy binary into result image from builder
COPY --from=builder /tmp/service .

# Copy grpc-health-probe into result image from builder
COPY --from=builder /go/bin/grpc-health-probe .

# Run binary
# We can specify additional options, such as:
#   --config=/path/to/config
#   --kube-config=/path/to/kubeconf
ENTRYPOINT ["/service"]
#CMD ["-logtostderr=true", "-v=1"]
