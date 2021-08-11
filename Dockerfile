FROM golang:1.16.5-buster as build

WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download -x

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./sender ./cmd/sender


FROM scratch

WORKDIR /app
COPY --from=build /app/sender /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/sender"]
