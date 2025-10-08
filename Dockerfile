FROM golang:1.25.1

WORKDIR /app

COPY go.mod go.sum/

RUN go mod download 

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o product-command ./cmd

EXPOSE 8080

ENTRYPOINT [ "./product-command" ]