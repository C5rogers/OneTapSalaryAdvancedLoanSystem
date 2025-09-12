FROM golang:latest AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a  -o main ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/main /app/

COPY ./data /app/data 



CMD ["./main"]
