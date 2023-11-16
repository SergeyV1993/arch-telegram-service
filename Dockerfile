FROM golang:1.21.3-alpine AS builder

RUN apk --no-cache add ca-certificates
RUN mkdir /telegram
ADD . /telegram
WORKDIR /telegram

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/srv ./

FROM scratch

COPY --from=builder /bin/srv /app/srv
COPY --from=builder /telegram/.env /.env
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/srv"]