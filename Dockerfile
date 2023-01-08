FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o typesense-exporter .

FROM busybox

WORKDIR /app

COPY --from=builder /app/typesense-exporter /usr/bin/

ENTRYPOINT ["typesense-exporter"]