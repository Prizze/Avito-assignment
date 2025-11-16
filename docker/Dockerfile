FROM golang:1.25.4-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
    && go install github.com/golang/mock/mockgen@latest

COPY . .

RUN go generate ./...

RUN go build -o /server ./cmd/app/main.go && chmod +x /server

EXPOSE 8080