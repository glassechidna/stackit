FROM golang:1.12-alpine AS deps
RUN apk add --update git upx curl
RUN curl -o gor.tgz -L https://github.com/goreleaser/goreleaser/releases/download/v0.106.0/goreleaser_Linux_x86_64.tar.gz
RUN tar -xvf gor.tgz && mv goreleaser /usr/bin

ENV CGO_ENABLED=0

WORKDIR /wd
COPY go.mod go.sum ./
RUN go mod download

FROM deps AS builder
COPY . .
RUN go build -ldflags="-s -w"
RUN upx stackit

FROM alpine
RUN apk add --update --no-cache ca-certificates
COPY --from=builder /wd/stackit /usr/bin
CMD ["/usr/bin/stackit"]
