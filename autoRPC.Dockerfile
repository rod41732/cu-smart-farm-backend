FROM intaniger/smartfarm_backend:latest

RUN rm -rf /go/src/github.com/rod41732/cu-smart-farm-backend
WORKDIR /go/src/github.com/rod41732/cu-smart-farm-backend
COPY . /go/src/github.com/rod41732/cu-smart-farm-backend
WORKDIR /go/src/github.com/rod41732/cu-smart-farm-backend/service
RUN go-wrapper download; exit 0
RUN go build main.go
RUN cp /go/src/github.com/rod41732/cu-smart-farm-backend/service/main /tmp/

FROM alpine
LABEL maintainer "Tanakorn Pisnupoomi"
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk --no-cache --update upgrade && apk --no-cache add ca-certificates
COPY --from=0 /tmp /backend
ENV GIN_MODE=release
WORKDIR /backend
CMD ["./main"]
EXPOSE 5555