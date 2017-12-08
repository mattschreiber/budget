FROM golang:1.8.1
ADD . /go/src/budget
RUN go get github.com/gorilla/mux
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/dgrijalva/jwt-go/request
RUN go get github.com/lib/pq
RUN go get golang.org/x/crypto/bcrypt
RUN go install budget
ENTRYPOINT /go/bin/budget
EXPOSE 5000
