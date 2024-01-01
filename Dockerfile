FROM golang:latest as builder

RUN update-ca-certificates

WORKDIR app/

COPY go.mod .

ENV GO111MODULE=on
RUN go mod download && go mod verify

COPY . .
# RUN apt-get update
# RUN apt-get -y install libdlib-dev libblas-dev libatlas-base-dev liblapack-dev libjpeg62-turbo-dev

RUN go build -o /app .

FROM debian:latest
# RUN apt-get update
# RUN apt-get -y install libdlib-dev libblas-dev libatlas-base-dev liblapack-dev libjpeg62-turbo-dev wget

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /usr/local/bin/app
# RUN mkdir /usr/local/bin/assets
# RUN wget -P /usr/local/bin/assets https://github.com/Kagami/go-face-testdata/raw/master/models/shape_predictor_5_face_landmarks.dat
# RUN wget -P /usr/local/bin/assets https://github.com/Kagami/go-face-testdata/raw/master/models/dlib_face_recognition_resnet_model_v1.dat
# RUN wget -P /usr/local/bin/assets https://github.com/Kagami/go-face-testdata/raw/master/models/mmod_human_face_detector.dat
EXPOSE 1000

ENTRYPOINT ["app"]
