FROM alpine:3.24.3@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a

RUN apk upgrade --no-cache \
    && apk add --no-cache bash ca-certificates curl procps tzdata vim \
    && cp /usr/share/zoneinfo/Asia/Taipei /etc/localtime \
    && echo "Asia/Taipei" > /etc/timezone

ENV TZ=Asia/Taipei

COPY _build/slctl /usr/local/bin
