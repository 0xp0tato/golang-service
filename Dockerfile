FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY shared/ ./shared/
COPY services/ ./services/
ARG SERVICE_NAME
RUN go build -o service ./services/${SERVICE_NAME}

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/service .
ARG PORT
EXPOSE ${PORT}
CMD ["./service"]