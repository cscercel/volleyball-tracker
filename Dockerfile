FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

ADD volleyball-tracker /usr/bin/volleyball-tracker

CMD ["volleyball-tracker"]
