FROM golang

RUN mkdir -p /go/src/github.com/ecosystemsoftware/ecosystem
WORKDIR /go/src/github.com/ecosystemsoftware/ecosystem
COPY . /go/src/github.com/ecosystemsoftware/ecosystem/
RUN go get -v