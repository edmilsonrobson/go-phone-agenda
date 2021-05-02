FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN ls
RUN go build -o main ./cmd/*.go

WORKDIR /dist

RUN cp /build/main .
COPY .env.docker .env

RUN chmod +x /dist/main
EXPOSE 8045

CMD ["/dist/main"]