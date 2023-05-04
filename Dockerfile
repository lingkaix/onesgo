FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN go mod download
ENV CGO_ENABLED=1
RUN go build -tags=nomsgpack -o app .

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/app /app/app
EXPOSE 8080
ENV JWT_KEY="JWT_secret-key!"
CMD ["./app"]

