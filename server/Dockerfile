# syntax=docker/dockerfile:1
# https://docs.docker.com/language/golang/build-images/

FROM golang:1.19-alpine

WORKDIR /app

COPY .env ./
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

COPY *.go ./

RUN go build -o /go-simple-live-polling

EXPOSE 8000

CMD [ "/go-simple-live-polling" ]