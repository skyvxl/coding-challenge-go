FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/app \
    ./cmd/main

RUN CGO_ENABLED=0 GOBIN=/out go install \
    github.com/pressly/goose/v3/cmd/goose@v3.27.3


FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/app /usr/local/bin/app
COPY --from=builder /out/goose /usr/local/bin/goose
COPY database/migrations ./migrations

USER app

EXPOSE 8090

CMD ["app"]