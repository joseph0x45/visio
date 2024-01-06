FROM golang:latest as builder
RUN update-ca-certificates
WORKDIR app/
COPY go.mod .
ENV GO111MODULE=on
RUN go mod download && go mod verify
COPY . .
RUN go build -o /app .
FROM debian:latest
# RUN apt-get update
# RUN apt-get -y install libdlib-dev libblas-dev libatlas-base-dev liblapack-dev libjpeg62-turbo-dev wget
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /usr/local/bin/app
COPY views /usr/local/bin/views
WORKDIR /usr/local/bin
EXPOSE 1000
ENTRYPOINT ["app"]
