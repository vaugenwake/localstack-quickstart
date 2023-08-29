# Localstack Quickstart

## WIP - Education Project
This project may become an open project to be used but at this stage it is purely an educational project to learn GoLang.

## Overview
At a high level I have used LocalStack in my local development environments and have concluded that a better tool is needed to manage setup and manage of it in docker environments.
This project is to build a new tool to allow you to provide a config file and execute a container on start-up of your project and have the yaml file read and executed to re-establish your ideal environment

### Sample config file
```YAML
connection:
  protocol: http
  endpoint: localstack
  port: 4566

resources:
  my-bucket:
    type: s3
    options:
      name: my-app-bucket
  my-queue:
    type: sqs
    options:
      name: queue-name
      dead_letter: my-queue-dlq
      depends_on:
        - my-queue-dlq
  my-queue-dlq:
    type: sqs
    options:
      name: dlq-name
```

### Adding to your docker environment
```YAML
version: '3'

networks:
  localstack:

services:
  app:
    image: vaugenwake/localstack-quickstart
    container_name: localstack_quickstart
    working_dir: /app
    environment:
      - AWS_ACCESS_KEY_ID=localhost
      - AWS_SECRET_ACCESS_KEY=localhost
      - AWS_REGION=us-east-1
    volumes:
      - ./config.yml:/app/config.yml # Mount your config file
    entrypoint: ["/localstack-quickstart", "-config=config.yml"]
    networks:
      - localstack
  localstack:
    container_name: "${LOCALSTACK_DOCKER_NAME-localstack_main}"
    image: localstack/localstack
    ports:
      - "127.0.0.1:4566:4566"            # LocalStack Gateway
      - "127.0.0.1:4510-4559:4510-4559"  # external services port range
    environment:
      - DEBUG=${DEBUG-}
      - DOCKER_HOST=unix:///var/run/docker.sock
    volumes:
      - "${LOCALSTACK_VOLUME_DIR:-./volume}:/var/lib/localstack"
      - "/var/run/docker.sock:/var/run/docker.sock"
    networks:
      - localstack
```

## Usage

### Initialize resources
```BASH
docker-compose run --rm app
```

## Todo:
* Get it to work with S3 buckets