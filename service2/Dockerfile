FROM golang

ADD . /workdir

WORKDIR /workdir

RUN go build -o app2 app.go   

EXPOSE 8081

CMD [ "./app2" ]