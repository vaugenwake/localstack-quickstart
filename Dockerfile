FROM golang:1.19.2-alpine3.16

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /localstack-quickstart

ENTRYPOINT [ "/localstack-quickstart" ]