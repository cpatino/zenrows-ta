FROM golang:1.25.7-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Set environment variables for MongoDB connection (use a secure method for production)
ENV MONGO_HOST=mongodb:27017
ENV MONGO_USER=user
ENV MONGO_PASSWORD=pass
ENV MONGO_DB=zenrows-ta-db

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]