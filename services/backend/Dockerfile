FROM golang:1.22-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum .
RUN go mod tidy
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /server

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /server server
ENTRYPOINT ["/app/server"]
