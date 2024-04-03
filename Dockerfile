FROM golang AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN swag init -g cmd/main.go -o ./docs
RUN CGO_ENABLED=0 GOOS=linux go build -o /main github.com/avialog/backend/cmd

FROM alpine
WORKDIR /
COPY --from=builder /main /main
EXPOSE 3000
ENTRYPOINT ["/main"]