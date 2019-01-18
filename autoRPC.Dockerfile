FROM intaniger/smartfarm_backend:latest

RUN rm -rf /go/src/github.com/rod41732/cu-smart-farm-backend
WORKDIR /go/src/github.com/rod41732/cu-smart-farm-backend
COPY . /go/src/github.com/rod41732/cu-smart-farm-backend
WORKDIR /go/src/github.com/rod41732/cu-smart-farm-backend/service
RUN go-wrapper download; exit 0
RUN go build main.go
CMD [ "./main" ]
EXPOSE 5555