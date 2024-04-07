FROM golang AS builder
WORKDIR /app
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN swag init --parseDependency --parseInternal -g cmd/main.go -o  ./docs
RUN CGO_ENABLED=0 GOOS=linux go build -o /main github.com/avialog/backend/cmd

FROM alpine
WORKDIR /
COPY --from=builder /main /main
EXPOSE 3000
ENTRYPOINT ["/main"]