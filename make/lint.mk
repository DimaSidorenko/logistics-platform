CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
LINTVER=v1.62.2
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}


bindir:
	mkdir -p ${BINDIR}


install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})


lint-cart:
	@if [ -f "cart/go.mod" ]; then \
		output=$$(${LINTBIN} --config=.golangci.yaml run cart 2>&1); \
		exit_code=$$?; \
		echo "$$output"; \
		if [ $$exit_code -ne 0 ]; then \
			if echo "$$output" | grep -q "no go files to analyze"; then \
				exit 0; \
			else \
				exit $$exit_code; \
			fi \
		fi \
	fi

lint-loms:
	@if [ -f "loms/go.mod" ]; then \
		output=$$(${LINTBIN} --config=.golangci.yaml run loms 2>&1); \
		exit_code=$$?; \
		echo "$$output"; \
		if [ $$exit_code -ne 0 ]; then \
			if echo "$$output" | grep -q "no go files to analyze"; then \
				exit 0; \
			else \
				exit $$exit_code; \
			fi \
		fi \
	fi

lint-notifier:
	@if [ -f "notifier/go.mod" ]; then \
		output=$$(${LINTBIN} --config=.golangci.yaml run notifier 2>&1); \
		exit_code=$$?; \
		echo "$$output"; \
		if [ $$exit_code -ne 0 ]; then \
			if echo "$$output" | grep -q "no go files to analyze"; then \
				exit 0; \
			else \
				exit $$exit_code; \
			fi \
		fi \
	fi

lint-comments:
	@if [ -f "comments/go.mod" ]; then \
		output=$$(${LINTBIN} --config=.golangci.yaml run comments 2>&1); \
		exit_code=$$?; \
		echo "$$output"; \
		if [ $$exit_code -ne 0 ]; then \
			if echo "$$output" | grep -q "no go files to analyze"; then \
				exit 0; \
			else \
				exit $$exit_code; \
			fi \
		fi \
	fi
