FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./server .


FROM alpine
WORKDIR /app
COPY --from=builder /app/server ./server
COPY db/migration ./db/migration
COPY start.sh .

EXPOSE 8080
CMD [ "/app/server" ]
ENTRYPOINT [ "/app/start.sh" ]