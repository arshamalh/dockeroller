FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
# Pure Go (SQLite included), so the binary is static.
RUN CGO_ENABLED=0 GOOS=linux go build -o /dockeroller .

FROM alpine:3.22 AS production
COPY --from=builder /dockeroller /usr/local/bin/dockeroller
# Sessions live in ./dockeroller.db; mount a volume here to keep them.
WORKDIR /data
ENTRYPOINT ["dockeroller"]
