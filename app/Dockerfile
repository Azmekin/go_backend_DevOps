FROM debian:latest
LABEL authors="Dmitry"
WORKDIR /backend
#COPY /backend .
RUN apt update && apt install golang-go -y

RUN  go version
RUN apt-get install openssl
RUN apt-get install ca-certificates
ARG cert_location=/usr/local/share/ca-certificates

# Get certificate from "github.com"

RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt
# Get certificate from "proxy.golang.org"
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt
# Update certificates
RUN update-ca-certificates


COPY /backend .
RUN go  mod init github.com/backend
RUN go get github.com/redis/go-redis/v9
RUN go get github.com/joho/godotenv
RUN go mod tidy

RUN  CGO_ENABLED=0 GOOS=linux go build -o /entrypoint
EXPOSE 8080
CMD ["/entrypoint"]
