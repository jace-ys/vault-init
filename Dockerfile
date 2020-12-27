FROM golang:1.15 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -a -installsuffix cgo ./cmd/vault-init/...

FROM alpine:3.11
COPY --from=builder /go/bin/vault-init /usr/local/bin/vault-init
ENTRYPOINT ["vault-init"]
