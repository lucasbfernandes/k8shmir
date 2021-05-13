FROM golang:1.14-alpine

RUN apk add --no-cache git

WORKDIR /app/req-smr

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/req-smr cmd/req-smr/*.go

EXPOSE 8080

ENTRYPOINT ["./out/req-smr"]
