FROM golang:1.13
LABEL contacts="github.com/peiiiajikuh" \
maintainer="peiiiajikuh"

WORKDIR /go/src/ascii-art-web
COPY . .

EXPOSE 8080

RUN go build -o main .

CMD ["./main"]