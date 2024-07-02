FROM golang:1.20-alpine

RUN apk update && apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/main .

RUN ls -al /app

EXPOSE 8080

ENTRYPOINT ["/app/main"]
