FROM golang:1.22-bookworm as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
ENV CGO_ENABLED=1 GO111MODULE=on GOPROXY=https://goproxy.cn
RUN go mod tidy && go mod download
RUN go test ./...
RUN cd cmd && \
    GIT_COMMIT=$(git rev-list -1 HEAD) && \
    BUILD_TIME=$(date "+%Y-%m-%d_%H:%M:%S") && \
    go build -ldflags "-X main.GitCommit=$GIT_COMMIT -X main.BuildTime=$BUILD_TIME" \
    -o /build/websockify .

FROM debian:bookworm
LABEL maintainer="shenjindi@ruzhila.cn"
RUN apt-get update && apt-get install -y ca-certificates tzdata
ENV DEBIAN_FRONTEND noninteractive
ENV LANG C.UTF-8

COPY --from=builder /build/websockify /usr/bin

EXPOSE 8000
ENTRYPOINT ["/usr/bin/websockify"]