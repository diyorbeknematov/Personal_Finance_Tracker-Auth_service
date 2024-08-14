FROM golang:1.22.2 AS builder

WORKDIR /auth-service

COPY . .
RUN go mod download
COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o ./../auth .

FROM alpine:latest

WORKDIR /auth-service

COPY --from=builder /auth-service/auth .
COPY --from=builder /auth-service/pkg/logs/app.log ./pkg/logs/
COPY --from=builder /auth-service/.env .

EXPOSE 4444

CMD [ "./auth" ]