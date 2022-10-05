FROM golang:1.19-alpine as builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /server server.go
FROM alpine:3
COPY --from=builder server /bin/server
ENTRYPOINT ["/bin/server"]