FROM golang:1.10

ADD . /go/src/github.com/niclasgeiger/crd-controller

WORKDIR /go/src/github.com/niclasgeiger/crd-controller

EXPOSE 80

ENTRYPOINT ["go", "run", "main.go"]