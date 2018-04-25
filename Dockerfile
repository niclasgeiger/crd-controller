FROM golang:1.10

ADD . /go/src/github.com/niclasgeiger/crd-controller

WORKDIR /go/src/github.com/niclasgeiger/crd-controller

ENTRYPOINT ["go", "run", "main.go"]