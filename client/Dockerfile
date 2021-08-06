FROM golang:1.16-alpine
ENV CGO_ENABLED=0

RUN apk add --no-cache git python3 build-base gcc wget openssh openssh-client curl bind-tools

WORKDIR /usr/app
COPY . .
RUN go mod download && go build main.go

EXPOSE 80 443

ENTRYPOINT [ "./main" ]
