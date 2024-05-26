FROM golang:alpine AS builder

WORKDIR /var/app

COPY . .

RUN go build cmd/rate-limiter/main.go

FROM scratch

WORKDIR /var/app

COPY --from=builder /var/app/main .
COPY --from=builder /var/app/dev.env .
COPY --from=builder /var/app/assets/tokens.json ./assets/

EXPOSE 8080

ENTRYPOINT [ "./main" ]