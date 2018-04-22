FROM golang:1.10-alpine3.7

ARG LIBRESSL_VERSION=2.7
ARG LIBRDKAFKA_VERSION=0.11.4-r1

RUN sed -i -e 's/v[[:digit:]]\.[[:digit:]]/edge/g' /etc/apk/repositories
RUN apk update && apk add --upgrade apk-tools && \
    apk add libressl${LIBRESSL_VERSION}-libcrypto libressl${LIBRESSL_VERSION}-libssl --update-cache
RUN apk add librdkafka=${LIBRDKAFKA_VERSION} --update-cache && \
    apk add librdkafka-dev=${LIBRDKAFKA_VERSION} --update-cache
RUN apk add git openssh openssl yajl-dev zlib-dev cyrus-sasl-dev openssl-dev build-base coreutils

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...



FROM alpine:edge

ARG LIBRESSL_VERSION=2.7
ARG LIBRDKAFKA_VERSION=0.11.4-r1

RUN apk update && apk add --upgrade apk-tools && \
    apk add libressl${LIBRESSL_VERSION}-libcrypto libressl${LIBRESSL_VERSION}-libssl --update-cache
RUN apk add librdkafka=${LIBRDKAFKA_VERSION} --update-cache && \
    apk add librdkafka-dev=${LIBRDKAFKA_VERSION} --update-cache
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /root/
COPY --from=0 /go/bin/app .

ENTRYPOINT ["/root/app"]