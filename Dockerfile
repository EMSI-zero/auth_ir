FROM golang:1.22.3-alpine3.20 as build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk add --no-cache make git

WORKDIR /go/src/github.com/emsi-zero/auth_ir

# Pulling dependencies
COPY ./Makefile ./go.* ./
RUN make deps

# Building stuff
COPY . /go/src/github.com/emsi-zero/auth_ir

# Make sure you change the RELEASE_VERSION value before publishing an image.
RUN RELEASE_VERSION=v0.0.1 make build

FROM alpine:3.20
RUN adduser -D -u 1000 myshop

RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/emsi-zero/auth_ir/auth_ir /usr/local/bin/auth_ir
COPY --from=build /go/src/github.com/emsi-zero/auth_ir/migrations /usr/local/etc/auth_ir/migrations/
RUN ln -s /usr/local/bin/auth_ir /usr/local/bin/gotrue

ENV GOTRUE_DB_MIGRATIONS_PATH /usr/local/etc/auth_ir/migrations

USER myshop
CMD ["auth_ir"]

