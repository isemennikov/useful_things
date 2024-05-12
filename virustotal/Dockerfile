# syntax=docker/dockerfile:1.4

#  set vars for build
ARG GO_VERSION=1.22.3
ARG OS_VERSION=alpine

# use base image for build 
FROM golang:${GO_VERSION}-${OS_VERSION} AS base

# git & make install 
RUN apk update && apk add --no-cache git make


# clone vt repo and install app
RUN git clone https://github.com/VirusTotal/vt-cli.git /vt-cli \
    && cd /vt-cli \
    && make install


ENV GOBIN=/root/go/bin
ENV PATH=$PATH:$GOBIN

# create finish image wich contains base alpine image
FROM alpine:latest

#copy exfiles 
COPY --from=base /go/bin/vt /usr/local/bin/vt

#set workdirictory
WORKDIR /data


#use secret while building 
RUN --mount=type=secret,id=apikey \
    API_KEY=$(cat /run/secrets/apikey) && \
    echo $API_KEY | vt init 

# set entrypoint
ENTRYPOINT [ "vt" ]

# default list --help


CMD ["help"]
