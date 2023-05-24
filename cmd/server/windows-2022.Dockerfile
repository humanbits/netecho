FROM golang:1.20-windowsservercore-ltsc2022 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


COPY cmd/server cmd/server

RUN CGO_ENABLED=0 GOOS=linux go build -o /netecho-server ./cmd/server/main.go


FROM scratch

COPY --from=builder /netecho-server /netecho-server

ENTRYPOINT ["/netecho-server"]