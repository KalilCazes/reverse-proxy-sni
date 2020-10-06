FROM golang

ADD . /workdir

WORKDIR /workdir

RUN go build -o app1 app.go   

EXPOSE 8080

CMD [ "./app1" ]