FROM golang:1.8.1

RUN go get gopkg.in/kataras/iris.v6
RUN go get github.com/jinzhu/gorm
RUN go get github.com/tarantool/go-tarantool
RUN go get gopkg.in/vmihailenco/msgpack.v2
RUN go get github.com/lib/pq
RUN go get github.com/markbates/goth
RUN go get github.com/gorilla/sessions
RUN go get github.com/centrifugal/gocent
RUN go get github.com/shopspring/decimal
RUN go get github.com/FireGM/speechkit
RUN go get github.com/FireGM/chats


ADD . /go/src/donategold.me
RUN go install donategold.me
ENTRYPOINT /go/bin/donategold.me

