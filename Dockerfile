FROM golang:1.12-alpine
RUN apk add --update git py2-pip
RUN pip install awscli

WORKDIR /wd
COPY go.mod go.sum ./
RUN go mod download
