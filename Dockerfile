FROM golang:1.21.0-alpine3.17

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /localstack-quickstart

ENTRYPOINT [ "/localstack-quickstart" ]