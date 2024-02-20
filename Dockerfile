FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./server .

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz


FROM alpine
WORKDIR /app
# COPY .env .env
COPY --from=builder /app/server ./server
COPY --from=builder /app/migrate ./migrate
COPY db/migration /app/migration
COPY start.sh .

EXPOSE 8080
CMD [ "/app/server" ]
ENTRYPOINT [ "/app/start.sh" ]