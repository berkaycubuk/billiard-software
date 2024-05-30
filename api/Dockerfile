FROM golang:1.20-alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

# COPY . .

# RUN go mod tidy

# COPY *.go ./

# RUN go build -o /billiard_software_api

CMD ["air", "-c", ".air.toml"]