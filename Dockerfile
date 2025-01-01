ARG GOLANG_VERSION=1.23.2

FROM golang:${GOLANG_VERSION}-alpine AS builder
ENV CGO_ENABLED 0

WORKDIR /usr/local/src

RUN apk --no-cache add bash

# dependencies
COPY ["./go.mod", "./go.sum", "./"]
RUN go mod download

# build
COPY ./ ./
RUN go build -o ./bin/app cmd/api/main.go

FROM alpine

COPY --from=builder /usr/local/src/bin/app /

# Копируем папку configs and .env с конфигами в контейнер
COPY --from=builder /usr/local/src/configs /configs
COPY --from=builder /usr/local/src/.env /.env


CMD ["/app"]