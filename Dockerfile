FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN apk --no-cache add ca-certificates
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/server

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/config /config
COPY --from=builder /app/internal/db/migration /migration
COPY --from=builder /bin/app /app
ENV CONFIG_PATH=/config/config-compose.yaml
CMD ["/app"]