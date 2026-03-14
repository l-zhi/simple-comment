FROM golang:1.26 AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /app/backend
RUN go mod download

COPY backend/. .

ENV GOPROXY=https://goproxy.cn,direct
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:3

WORKDIR /app

COPY --from=builder /app/server /app/server

ENV APP_PORT=8081
ENV DB_HOST=mysql
ENV DB_PORT=3306
ENV DB_USER=root
ENV DB_PASSWORD=password
ENV DB_NAME=simple_comment

EXPOSE 8081

ENTRYPOINT ["/app/server"]

