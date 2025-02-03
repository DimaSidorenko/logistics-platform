.PHONY: build
build:
	cd cart     && GOOS=linux GOARCH=amd64 make build
	cd loms     && GOOS=linux GOARCH=amd64 make build
	cd notifier && GOOS=linux GOARCH=amd64 make build
	cd comments && GOOS=linux GOARCH=amd64 make build
