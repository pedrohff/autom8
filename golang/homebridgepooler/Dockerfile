# stage de build
FROM golang:1.20.1-alpine3.17 AS build

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY *.go /app
RUN go env

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o homebridgepooler main.go

# stage imagem final
FROM alpine:3.17

WORKDIR /app

COPY --from=build /app/homebridgepooler ./

EXPOSE 8000

CMD [ "./homebridgepooler" ]