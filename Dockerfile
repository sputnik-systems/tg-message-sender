FROM golang:1.20 as build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./sender ./cmd/sender


FROM scratch

WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/sender /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/sender"]
