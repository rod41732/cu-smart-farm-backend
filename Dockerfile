FROM golang:1.9

WORKDIR /go/src/app
COPY . /go/src/app
RUN go-wrapper download; exit 0
RUN go build main.go