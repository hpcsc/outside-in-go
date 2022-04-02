FROM golang:1.18.0 as builder
WORKDIR /build

ADD . .

RUN GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go build -v --ldflags '-extldflags -static' \
    -o outside-in-go \
    ./cmd/outside-in-go/main.go

FROM scratch
WORKDIR /app
ENTRYPOINT ["/app/outside-in-go"]
COPY --from=builder /build/outside-in-go .
