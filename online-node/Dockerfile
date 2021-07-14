FROM golang:1.16-alpine
ENV CGO_ENABLED=0
ENV CONSUL_VERSION=1.10.0

RUN apk add --no-cache git python3 build-base gcc wget openssh openssh-client curl bind-tools
RUN curl -sSLo /tmp/consul.zip https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_amd64.zip && \
    unzip -d /bin /tmp/consul.zip && \
    rm -f /tmp/consul.zip && \
    mkdir -p /consul/config /consul/data /consul/ui 

WORKDIR /usr/app
COPY . .
RUN go mod download
RUN go build main.go
EXPOSE 3010

RUN chmod +x main
RUN chmod +x start.sh
CMD ["./start.sh"]