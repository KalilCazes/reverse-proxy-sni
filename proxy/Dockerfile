FROM golang

ADD . /workdir

WORKDIR /workdir

RUN go build -o proxy main.go   

EXPOSE 8000

CMD [ "./proxy" ]