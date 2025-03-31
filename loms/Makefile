BINDIR=${CURDIR}/bin
PACKAGE=route256/loms

# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN:=$(CURDIR)/bin
NOTES_PROTO_PATH:="protobuf/rpc/server"

bindir:
	mkdir -p ${BINDIR}

# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

# Устанавливаем proto описания validate
vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
		cd vendor-proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor-proto/validate
		mv vendor-proto/tmp/validate vendor-proto/
		rm -rf vendor-proto/tmp

# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

# Устанавливаем proto описания protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem

.PHONY: .vendor-rm
.vendor-rm:
	rm -rf vendor-proto

# Вендоринг внешних proto файлов
.PHONY: .vendor-proto
.vendor-proto: .vendor-rm  vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options


.PHONY: .bin-deps
.bin-deps: bindir
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4 && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1 && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1 && \
	GOBIN=$(LOCAL_BIN) go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.5 && \
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.1 && \
	GOBIN=$(LOCAL_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0


.PHONY: .protoc-generate
.protoc-generate:
	protoc \
	-I ${NOTES_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=${LOCAL_BIN}/protoc-gen-go \
	--go_out pkg/${NOTES_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=${LOCAL_BIN}/protoc-gen-go-grpc \
	--go-grpc_out pkg/${NOTES_PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=${LOCAL_BIN}/protoc-gen-grpc-gateway \
	--grpc-gateway_out pkg/${NOTES_PROTO_PATH} \
	--grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-validate=${LOCAL_BIN}/protoc-gen-validate \
	--validate_out="lang=go,paths=source_relative:pkg/protobuf/rpc/server" \
	${NOTES_PROTO_PATH}/service.proto
	go mod tidy


MIGRATION_NAME ?= add_registry

create-migration-file:
	$(LOCAL_BIN)/goose -dir migrations create -s $(MIGRATION_NAME) sql

up-migrations:
	$(LOCAL_BIN)/goose -dir migrations postgres "postgresql://user:password@127.0.0.1:5432/route256?sslmode=disable" up

revert-migration:
	$(LOCAL_BIN)/goose -dir migrations postgres "postgresql://user:password@127.0.0.1:5432/route256?sslmode=disable" down

compile-sql:
	$(LOCAL_BIN)/sqlc generate

build: .protoc-generate
	go run cmd/service/main.go

lint:
	golangci-lint run ./...

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

run-all: lint coverage build
	echo "Запуск всех задач завершен."

# Цель для запуска миграций
run-migrations-local: .bin-deps
	$(info Running migrations for loms...)
	$(LOCAL_BIN)/goose -dir migrations postgres "postgresql://loms-user:loms-password@localhost:5432/loms_db?sslmode=disable" up

# Цель для запуска миграций
run-migrations:
	$(info Running migrations for loms...)
	goose -dir ${MIGRATION_DIR} postgres "postgresql://loms-user:loms-password@postgres-master:5432/loms_db?sslmode=disable" up