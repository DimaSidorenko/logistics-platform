BINDIR=${CURDIR}/cmd
PACKAGE=route256/cart

bindir:
	mkdir -p ${BINDIR}

build: bindir
	#CONFIG_FILE=./configs/values_local.yaml go run cmd/main.go
	go run cmd/main.go
lint:
	golangci-lint run ./...

run-all: lint build
	echo "Запуск всех задач завершен."
