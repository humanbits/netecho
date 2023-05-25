FROM golang:1.20-buster as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


COPY cmd/client cmd/client

RUN CGO_ENABLED=0 GOOS=linux go build -o /netecho-client ./cmd/client/main.go


FROM scratch

COPY --from=builder /netecho-client /netecho-client

ENTRYPOINT ["/netecho-client"]