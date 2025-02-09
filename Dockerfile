# Build stage
FROM golang:1.17 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ring_buffer .

# Final stage
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /app/ring_buffer .
CMD ["./ring_buffer"]