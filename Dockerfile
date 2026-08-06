FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o proxify \
    ./cmd/proxify

FROM alpine:3.22

RUN addgroup -S proxify && \
    adduser -S proxify -G proxify

WORKDIR /app

COPY --from=builder /app/proxify .
COPY --from=builder /app/configs ./configs

USER proxify

EXPOSE 8080

ENTRYPOINT ["./proxify"]