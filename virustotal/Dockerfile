# syntax=docker/dockerfile:1.4

#  set vars for build
ARG GO_VERSION=1.22.3
ARG OS_VERSION=alpine.3.19

# use base image for build 
FROM golang:${GO_VERSION}-${OS_VERSION} AS base

# git & make install 
RUN apk update && apk add --no-cache git make

# clone vt repo and install app
RUN git clone https://github.com/VirusTotal/vt-cl.git /vt-cl \
    && /vt-cli \
    make install

# set var for go 
ENV GOBIN=/go/bin

# exfile vt copy to GOBIN
RUN cp /vt-cli/vt ${GOBIN}

