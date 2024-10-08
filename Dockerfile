FROM golang:1.23.2

WORKDIR /app

RUN apt-get update && apt-get install bash

ENV GO111MODULE on

# The command must have a path through, so install with GOMODULE off
COPY . /app

RUN go install github.com/air-verse/air@latest \
	&& go mod download && go mod tidy

RUN go build -o /tmp/main .

RUN echo 'PATH=$PATH:/app/bin' > /root/.bashrc

RUN curl -SL https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait -o /wait
RUN chmod +x /wait