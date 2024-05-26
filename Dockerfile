FROM golang:1.21 as build
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun ./cmd/rate-limiter

FROM scratch
WORKDIR /app
COPY --from=build /app/cloudrun .
COPY --from=build /app/dev.env .
COPY --from=build /app/resources/tokens.json ./resources/
EXPOSE 8080
ENTRYPOINT ["./cloudrun"]