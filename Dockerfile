# build
FROM golang:1.18-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN ls -la

RUN cd cmd/gophermart && go build -o /gophermart

# deploy
FROM alpine

WORKDIR /

COPY --from=build /gophermart /gophermart

EXPOSE 8080

ENTRYPOINT ["/gophermart"]




CMD ["./cmd/gophermart/gophermart"]
