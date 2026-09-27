# 多阶段构建
FROM golang:1.26 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /grpc-server ./server

FROM alpine:3.20
COPY --from=build /grpc-server /grpc-server
EXPOSE 50051
ENTRYPOINT ["/grpc-server"]
