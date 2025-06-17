# stage 1
FROM golang:1.23-alpine as builder
WORKDIR /app

RUN apk add --no-cache git libwebp-dev build-base

COPY go.* ./
RUN go mod download

# RUN go install github.com/air-verse/air@latest
COPY . ./
COPY .env.docker /app/.env
#RUN go build -o main cmd/tex/main.go
RUN GOOS=linux GOARCH=amd64 go build -o app ./cmd/tex

# stage 2
FROM alpine:3.19 as app
RUN apk --no-cache upgrade && apk --no-cache add ca-certificates

RUN apk add --no-cache bash postgresql-client ffmpeg libwebp

COPY --from=builder /app/app /usr/local/bin/app
COPY --from=builder /app/.env /usr/local/bin/.env
COPY --from=builder /app/scripts /usr/local/bin/scripts
COPY --from=builder /app/schemas /usr/local/bin/schemas

WORKDIR /usr/local/bin/

EXPOSE 7000
CMD ["/app"]
# CMD ["air", "-c", ".air.toml"]