# stage de build
FROM golang:1.20.1-alpine3.17 AS build

WORKDIR /app

COPY go.work /app/
COPY go.work.sum /app/

{{.SubmoduleModsList}}

RUN go mod download

{{.SubmoduleList}}


RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o homebridgepooler /app/homebridgepooler

# stage imagem final
FROM alpine:3.17

WORKDIR /app

COPY --from=build /app/homebridgepooler ./

EXPOSE 8000

CMD [ "./homebridgepooler" ]