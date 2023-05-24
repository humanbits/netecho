FROM golang:1.20-windowsservercore-ltsc2022 as builder

WORKDIR C:\\app

COPY go.mod go.sum ./

RUN go mod download


COPY cmd\server cmd\server

ENV CGO_ENABLED=0
RUN go build -o C:\\netecho-server .\cmd\server\main.go



COPY --from=builder C:\\netecho-server C:\\netecho-server

ENTRYPOINT ["C:\\netecho-server"]