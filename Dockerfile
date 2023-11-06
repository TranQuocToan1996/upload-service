FROM golang:1.21 AS builder

WORKDIR /app

COPY . /app/

RUN go mod tidy  && \
    make build cmd="api" service_name="upload" &&  \
    mkdir ./build  && \
    cp -p ./cmd/api/upload ./build && \
    cp -p ./cmd/api/.env ./build

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build /app
EXPOSE 8080

CMD ["./upload"]