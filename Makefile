setup:
	go install github.com/onsi/ginkgo/v2/ginkgo
	go install go.uber.org/mock/mockgen@latest
test:
	go generate ./...
	ginkgo -r -v ./...
gen:
	swag init -g cmd/main.go -o ./docs