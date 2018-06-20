FROM golang:1.10

MAINTAINER Niclas Geiger

ADD main /

ENTRYPOINT ["/main"]