TIMEOUT_UNIT = 3m

FORCE:

$(OUTPUT_DIR)/%: cmd/% FORCE
	go build -mod=vendor -v -o $@ ./$<

.PHONY: vendor
vendor:
	@echo Generating vendor directory
	@go mod tidy -compat=1.17 && go mod vendor


.PHONY: test
test: test-clean test-unit ## Run test-unit

test-clean:  ## Clean test cache
	@echo "Cleaning test cache"
	@go clean -testcache ./...

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	@go test -v -timeout $(TIMEOUT_UNIT)  ./... | { grep -v 'no test files'; true; }
