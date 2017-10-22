FROM golang:1.8.1

RUN go get gopkg.in/kataras/iris.v6
RUN go get github.com/centrifugal/gocent
RUN go get github.com/FireGM/chats
RUN go get github.com/jinzhu/gorm
RUN go get github.com/lib/pq


ADD . /go/src/donategold.me/chats
RUN go install donategold.me/chats
ENTRYPOINT /go/bin/chats