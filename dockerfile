FROM golang:1.19-alpine as builder
WORKDIR /build
COPY app/go.mod .
COPY app/go.sum .
RUN go mod download
COPY app .
RUN go build -o /server ./cmd/web/.
FROM alpine:3
COPY --from=builder server /bin/server