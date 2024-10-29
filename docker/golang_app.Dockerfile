FROM golang:alpine

RUN apk add --no-cache build-base git inotify-tools && \
    go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /golang_app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["sh", "/golang_app/docker/run.sh"]
