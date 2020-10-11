FROM golang

ENV GO111MODULE=on

WORKDIR /workdir

ADD . /workdir

RUN go mod download && go build -o proxy/proxy proxy/proxy.go   

EXPOSE 8000

CMD [ "proxy/proxy" ]