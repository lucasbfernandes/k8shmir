FROM golang:1.16-alpine

RUN apk add --no-cache git

WORKDIR /app/k8s-smr

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY ../../cmd/k8s-smr .

RUN go build -o ./out/k8s-smr cmd/k8s-smr/*.go

EXPOSE 8080

ENTRYPOINT ["./out/k8s-smr"]
