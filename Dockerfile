FROM golang:1.24.1-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY . .

RUN go mod download

RUN touch swiftcodes.db

RUN go build -o swift-code-app .

EXPOSE 8080

CMD ["./swift-code-app"]