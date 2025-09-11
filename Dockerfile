FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go


FROM alpine
WORKDIR /app

COPY --from=builder /app/main .
COPY ./db/migration ./db/migration
COPY ./app.env .
COPY ./start.sh .
COPY ./wait-for.sh .

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]