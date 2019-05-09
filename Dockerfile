FROM golang:1.12-alpine
RUN apk add --update git py2-pip upx
RUN pip install awscli
ADD https://github.com/goreleaser/goreleaser/releases/download/v0.106.0/goreleaser_Linux_x86_64.tar.gz .
RUN tar -xvf goreleaser_Linux_x86_64.tar.gz && mv goreleaser /usr/bin

WORKDIR /wd
COPY go.mod go.sum ./
RUN go mod download

ENV CGO_ENABLED=0
