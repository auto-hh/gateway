FROM golang:1.26.1-alpine3.23 AS builder

WORKDIR /gateway/

COPY ./go.mod ./

COPY ./cmd/ ./cmd/
COPY ./config/ ./config/
COPY ./internal/ ./internal/
COPY ./pkg/ ./pkg/

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/ ./...


FROM alpine:3.23 AS service

WORKDIR /gateway/

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /gateway/build/gateway ./

CMD ["sh", "-c", "./gateway"]
