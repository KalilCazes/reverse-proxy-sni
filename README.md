# reverse-proxy-sni
This is a simple example of a [reverse proxy](https://en.wikipedia.org/wiki/Reverse_proxy) that uses HTTPS and [SSL SNI](https://en.wikipedia.org/wiki/Server_Name_Indication) writen in go

## Running in a docker container
1. Generate the certificates running the [initKey.sh](https://github.com/KalilCazes/reverse-proxy-sni/blob/main/initKey.sh) script. This script also modifies /etc/hosts, so it needs to be **run as sudo**
2. Now you can use [docker-compose](https://github.com/docker/compose). To check out the reverse proxy, run **docker-compose up**. It'll create a docker image and start
3. Go to the browser and access https://localhost1:8000 and https://localhost2:8000
- By default the reverse proxy on docker listen on **port 8000** and the backend servers on **8080** and **8081**

## Running without docker
1. Repeat the first instruction above to generate the certificates
2. **go run proxy/proxy.go** to run the reverse proxy
3. **go run service1/app.go** and **go run service2/app.go** to run the backend services
4. Go to the browser and access https://localhost1:8000 and https://localhost2:8000

## Configuration file
This project uses [config.yaml](https://github.com/KalilCazes/reverse-proxy-sni/blob/main/config.yaml) as input to specify **port for reverse proxy and path to certificates**
## Running tests
1. **go test ./proxy -v**
