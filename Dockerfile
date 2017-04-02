FROM golang

RUN mkdir -p /go/src/github.com/jpincas/ghost
WORKDIR /go/src/github.com/jpincas/ghost
COPY . /go/src/github.com/jpincas/ghost/
RUN go get -v