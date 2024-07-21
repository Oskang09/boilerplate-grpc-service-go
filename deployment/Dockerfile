# about version selections
#
# Update on linux alpine os, debian os in case your library used latest alpine etc golang-1.20. 
# Mostly workaround would be update the docker version, and for golang will be 1 more step using bullseye build instead alpine
#
# Reference: https://github.com/docker-library/golang/issues/467#issuecomment-1601845758
#
FROM golang:1.20.7-bullseye as builder

ADD . /go/gateway

WORKDIR  /go/gateway

# RUN git config --global url."git@gitlab.revenuemonster.my:".insteadOf "https://gitlab.revenuemonster.my/"

RUN make

# Delete all go files
RUN find . -name "*.go" -type f -delete

FROM alpine:latest

COPY --from=builder /go/gateway /go/gateway

# ENV GOPRIVATE 'gitlab.revenuemonster.my/dinar-wallet/*'

WORKDIR  /go/gateway

ENTRYPOINT ./gateway
