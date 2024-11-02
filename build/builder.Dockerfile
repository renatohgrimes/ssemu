FROM golang:1.23
RUN apt upgrade -y && apt update -y && apt install -y mingw-w64
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download