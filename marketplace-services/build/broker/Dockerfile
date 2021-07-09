FROM i386/ubuntu:latest

LABEL maintainer="michael.sober@ymail.com"
LABEL version="1.0"
LABEL description="Dockerfile for the angular frontend of the data marketplace proxy"

WORKDIR /app

COPY ./bin/broker .

CMD ["./broker"]