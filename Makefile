SERVER_APP_NAME=fastkey-server
CLI_APP_NAME=fastkey-cli

.PHONY: build-server
build-server:
	go build -o ${SERVER_APP_NAME} cmd/server/main.go

.PHONY: build-cli
build-cli:
	go build -o ${CLI_APP_NAME} cmd/cli/main.go

.PHONY: run-server
run-server: build-server
	./${SERVER_APP_NAME}

.PHONY: run-server-with-config
run-server-with-config: build-server
	CONFIG_FILE_NAME=config.yaml ./${SERVER_APP_NAME}

.PHONY: run-cli
run-cli: build-cli
	./${CLI_APP_NAME} $(ARGS)

.PHONY: run_unit_test
run_unit_test:
	go test ./internal/...

.PHONY: run_e2e_test
run_e2e_test: build-server
	go test ./test/...

.PHONY: run_test_coverage
run_test_coverage:
	go test ./... -coverprofile=coverage.out