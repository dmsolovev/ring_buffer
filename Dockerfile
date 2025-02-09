# Build stage
FROM golang:latest AS build
RUN mkdir -p /go/src/ring_buffer
WORKDIR /go/src/ring_buffer
ADD main.go .
ADD go.mod .
RUN go install .

# Final stage
FROM alpine:latest 
LABEL version="1.0.0" 
LABEL maintainer="DmSolovyov<dmsolovyov@gmail.com>"
WORKDIR /root/
COPY --from=build /go/bin/ring_buffer .
ENTRYPOINT [./ring_buffer]
