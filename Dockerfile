FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -o dclr .

FROM alpine:3.18 AS production
COPY --from=builder /app/dclr .
ENTRYPOINT [ "./dclr" ]
