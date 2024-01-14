FROM golang:latest

RUN apt-get update -y

RUN apt-get install -y libdlib-dev libblas-dev libatlas-base-dev liblapack-dev libjpeg62-turbo-dev

RUN go install github.com/cosmtrek/air@latest
