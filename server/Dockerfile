FROM golang:1.16-alpine
ENV CGO_ENABLED=0
# ENV GODEBUG=http2debug=2

RUN apk add --no-cache git python3 build-base gcc wget openssh openssh-client curl bind-tools

WORKDIR /usr/app
COPY . .
RUN go mod download
RUN go build main.go && chmod +x main

EXPOSE 3010

ENTRYPOINT [ "./main" ]
