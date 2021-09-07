FROM  golang:1.17.0-alpine3.13 AS builder

WORKDIR /code
COPY go.mod .
COPY go.sum .
RUN go mod download
ADD . .

# build the binary
RUN env CGO_ENABLED=0 GOOS=linux go build -o /monitor cmd/monitor.go

# final stage
FROM alpine:3.14.2

COPY --from=builder /monitor /

RUN chmod +x /monitor