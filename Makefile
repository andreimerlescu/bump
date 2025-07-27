# Go Application Makefile by github.com/andreimerlescu (Apache 2.0 License)
#
MAIN_PATH := $(shell realpath .)
APP_NAME := $(shell basename "$(shell realpath $(MAIN_PATH))")
PACKAGE_PATH := $(shell realpath "./$(APP_NAME)/")
BIN_DIR=bin
TEST_DIR=test-results
APP_BINARY := $(BIN_DIR)/$(APP_NAME)
MAKE := make
ECHO := echo
GO := go

# Go build flags
# -s: Strip symbols (reduces binary size)
# -w: Omit DWARF debugging information
LDFLAGS=-ldflags "-s -w"

export APP_BINARY
export MAIN_PATH
export APP_NAME

$(TEST_DIR):
	@if [ ! -d $(TEST_DIR) ]; then \
		mkdir -p $(TEST_DIR); \
	fi

$(BIN_DIR):
	@if [ ! -d $(BIN_DIR) ]; then \
		mkdir -p $(BIN_DIR); \
	fi

.PHONY: all clean summary install app-binary

all: summary darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 windows-amd64 install

clean:
	@rm -rf $(BIN_DIR)
	@echo "REMOVED: $(BIN_DIR)"

summary:
	@if ! command -v summarize > /dev/null; then \
		go install github.com/andreimerlescu/summarize@latest; \
	fi
	@summarize -i "go,Makefile,md,mod,sum,LICENSE,gitignore"

app-binary: $(BIN_DIR)
	@printf "Building binary target: %s/%s\n" "${GOOS}" "${GOARCH}"
	@if [ -f "$(APP_BINARY)" ]; then \
		rm -f "$(APP_BINARY)"; \
	fi
	@if [ "${GOOS}" == "windows" ]; then \
		$(GO) build $(LDFLAGS) -o "$(APP_BINARY).exe" $(MAIN_PATH); \
	else \
		$(GO) build $(LDFLAGS) -o "$(APP_BINARY)-${GOOS}-${GOARCH}" $(MAIN_PATH); \
	fi

install: $(BIN_DIR)
	@if [[ "$(shell go env GOOS)" == "windows" ]]; then \
		cp "$(BIN_DIR)/$(APP_NAME).exe" "$(shell go env GOBIN)/$(APP_NAME).exe"; \
	else \
		cp "$(BIN_DIR)/$(APP_NAME)-$(shell go env GOOS)-$(shell go env GOARCH)" "$(shell go env GOBIN)/$(APP_NAME)"; \
	fi
	@echo "NEW: $(shell which $(APP_NAME))"

.PHONY: darwin-amd64 darwin-amd64 linux-amd64 linux-arm64 windows-amd64

darwin-amd64: $(BIN_DIR) summary
	@GOOS=darwin GOARCH=amd64 $(MAKE) app-binary

darwin-arm64: $(BIN_DIR) summary
	@GOOS=darwin GOARCH=arm64 $(MAKE) app-binary

linux-arm64: $(BIN_DIR) summary
	@GOOS=linux GOARCH=arm64 $(MAKE) app-binary

linux-amd64: $(BIN_DIR) summary
	@GOOS=linux GOARCH=amd64 $(MAKE) app-binary

windows-amd64: $(BIN_DIR) summary
	@GOOS=windows GOARCH=amd64 $(MAKE) app-binary

.PHONY: test test-unit test-fuzz test-bench

UNIT_LOG="$(MAIN_PATH)/$(TEST_DIR)/results.unit.md"
FUZZ_LOG="$(MAIN_PATH)/$(TEST_DIR)/results.fuzz.md"
BENCH_LOG="$(MAIN_PATH)/$(TEST_DIR)/results.benchmark.md"

test: test-unit test-bench test-fuzz

test-unit: $(TEST_DIR)
	@printf "%s" "Testing Unit... "
	@touch $(UNIT_LOG)
	@echo "## \`$(UNIT_LOG)\` \n\n Test results captured at $(shell date +"%Y-%m-%d %H:%M:%S"). \n\n\`\`\`log" > $(UNIT_LOG)
	@start_time=$$(date +%s); \
	cd $(PACKAGE_PATH) && go test -vet=all -count=1 ./... >> $(UNIT_LOG); \
	test_result=$$?; \
	end_time=$$(date +%s); \
	elapsed=$$((end_time - start_time)); \
	if [ $$test_result -eq 1 ]; then \
		echo "FAILED!"; \
		exit 1; \
	fi; \
	echo "\`\`\`" >> $(UNIT_LOG); \
	echo "" >> $(UNIT_LOG); \
	echo "SUCCESS! Took $$elapsed (s)! Wrote $(shell basename "$(UNIT_LOG)") ( size: $(shell du -h "$(UNIT_LOG)" | awk '{print $$1}') )"

test-fuzz: $(TEST_DIR)
	@printf "%s" "Testing Fuzz... "
	@touch $(FUZZ_LOG)
	@echo "## \`$(FUZZ_LOG)\` \n\n Test results captured at $(shell date +"%Y-%m-%d %H:%M:%S"). \n\n\`\`\`log" > $(FUZZ_LOG)
	@start_time=$$(date +%s); \
	cd $(PACKAGE_PATH) && go test -vet=all -count=1 -fuzz=Fuzz -fuzztime=31s >> $(FUZZ_LOG); \
	test_result=$$?; \
	end_time=$$(date +%s); \
	elapsed=$$((end_time - start_time)); \
	if [ $$test_result -eq 1 ]; then \
		echo "FAILED!"; \
		exit 1; \
	fi; \
	echo "\`\`\`" >> $(FUZZ_LOG); \
	echo "" >> $(FUZZ_LOG); \
	echo "SUCCESS! Took $$elapsed (s)! Wrote $(shell basename "$(FUZZ_LOG)") ( size: $(shell du -h "$(FUZZ_LOG)" | awk '{print $$1}') )"

test-bench: $(TEST_DIR)
	@printf "%s" "Testing Benchmark... "
	@echo "## \`$(BENCH_LOG)\` \n\n Test results captured at $(shell date +"%Y-%m-%d %H:%M:%S"). \n\n\`\`\`log" > $(BENCH_LOG)
	@start_time=$$(date +%s); \
	cd $(PACKAGE_PATH) && go test -vet=all -count=1 -bench=.  >> $(BENCH_LOG); \
	test_result=$$?; \
	end_time=$$(date +%s); \
	elapsed=$$((end_time - start_time)); \
	if [ $$test_result -eq 1 ]; then \
		echo "FAILED!"; \
		exit 1; \
	fi; \
	echo "\`\`\`" >> $(BENCH_LOG); \
	echo "" >> $(BENCH_LOG); \
	echo "SUCCESS! Took $$elapsed (s)! Wrote $(shell basename "$(BENCH_LOG)") ( size: $(shell du -h "$(BENCH_LOG)" | awk '{print $$1}') )"
