FROM golang:1.23.4-alpine As builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN go build -o main .

FROM gcr.io/distroless/base-debian10

COPY --from=builder /app/main /app/main
COPY ./docs /app/docs

# DIhapus jika .envnya tidak ingin dijalankan
COPY .env /app/.env


WORKDIR /app

EXPOSE 8000

CMD ["/app/main"]
