BINDIR=${CURDIR}/bin
PACKAGE=route256/notifier

bindir:
	mkdir -p ${BINDIR}

lint:
	golangci-lint run ./...

build: bindir
	go build -o ${BINDIR}/consumer_group ./cmd/main/main.go

run-all: lint build
	${BINDIR}/consumer_group & \
	${BINDIR}/consumer_group & \
	${BINDIR}/consumer_group
	echo "Все инстансы consumer_group запущены."
