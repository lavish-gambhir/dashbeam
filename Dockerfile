
FROM golang:1.23.3-alpine AS builder
LABEL maintainer="lavish_gambhir@icloud.com"

RUN apk update && apk add --no-cache build-base

WORKDIR /app

COPY go.work go.work.sum ./
COPY go.mod go.sum ./

COPY services/ ./services/
COPY shared/ ./shared/
COPY pkg/ ./pkg/
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server ./cmd/server

# Build migrator separately
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/migrator ./cmd/migrator

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/bin/server ./bin/server
COPY --from=builder /app/bin/migrator ./bin/migrator
COPY configs/ ./configs/
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8080

# ENV APP_ENV=staging // TODO

CMD ["./bin/server"]
