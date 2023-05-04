FROM golang:1.20
# // ? use multi-stage build to have a smaller and cleaer image
WORKDIR /app
COPY . .
RUN go mod download
ENV CGO_ENABLED=1
RUN go build -tags=nomsgpack -o app .
EXPOSE 8080
ENV JWT_KEY="JWT_secret-key!"
CMD ["./app"]
