FROM golang:1.9

# RUN mkdir /go/src/github.com
WORKDIR /go/src/github.com/rod41732/cu-smart-farm-backend
COPY . /go/src/github.com/rod41732/cu-smart-farm-backend
RUN go-wrapper download; exit 0
RUN go build main.go
EXPOSE 3000