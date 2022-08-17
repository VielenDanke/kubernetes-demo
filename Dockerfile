FROM golang:latest AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o /books

FROM alpine:latest

WORKDIR /

COPY --from=build /books /books

EXPOSE 9090

ENTRYPOINT ["./books"]