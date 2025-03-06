BINDIR=${CURDIR}/cmd
PACKAGE=route256/cart

bindir:
	mkdir -p ${BINDIR}

build: bindir
	#CONFIG_FILE=./configs/values_local.yaml go run cmd/main.go
	go run cmd/main.go
lint:
	golangci-lint run ./...

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

run-all: lint coverage build
	echo "Запуск всех задач завершен."
