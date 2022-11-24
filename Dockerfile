FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /stackup-bundler

EXPOSE 4337

CMD ["/stackup-bundler", "start", "--mode", "private"]
