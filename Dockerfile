FROM golang:1.12
RUN apt-get update && apt-get install -y awscli

WORKDIR /wd
COPY go.mod go.sum ./
RUN go mod download
