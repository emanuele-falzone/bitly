generate-code:
	# Install mockgen
	go install github.com/golang/mock/mockgen@v1.6.0
	
	# Generate mocks for test purposes
	# Generate mock redirection repository
	mockgen -destination=./test/mock/redirection_repository.go \
		-package=mock \
		-mock_names=Repository=MockRedirectionRepository \
		github.com/emanuelefalzone/bitly/internal/application/redirection Repository

	# Generate mock event repository
	mockgen -destination=./test/mock/event_repository.go \
		-package=mock \
		-mock_names=Repository=MockEventRepository \
		github.com/emanuelefalzone/bitly/internal/application/event Repository

	# Generate mock key generator
	mockgen -destination=./test/mock/key_generator_service.go \
		-package=mock \
		github.com/emanuelefalzone/bitly/internal/application/service KeyGenerator

	# Install protobuf and grpc
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/protobuf/cmd/protoc-gen-go

	# Generate protobuf and grpc server
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		internal/port/grpc/pb/bitly_service.proto

generate-docs:
	# Ensure docs directory exists
	mkdir -p ./docs

	# Install proto documentation generator
	go get github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v1.5.1
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

	# Generate documentation for grpc service
	protoc \
		--doc_out=./docs \
  		--doc_opt=html,proto.html \
		internal/port/grpc/pb/bitly_service.proto

	# Install swag
	go get github.com/swaggo/swag/gen@v1.8.2
	go get github.com/swaggo/swag/cmd/swag@v1.8.2
	go install github.com/swaggo/swag/cmd/swag

	# Generate documentation for grpc service
	swag init -d internal/port/http --generalInfo server.go

build-for-production:
	# Build removing debug info
	go build -ldflags "-s -w" -v ./cmd/main.go

build-for-development:
	# Build using -race to check for race conditions while running
	go build -race -v ./cmd/main.go

CVPKG=$(go list ./internal/... | grep -v pb | tr '\n' ',')

run-unit-tests:
	# Run unit tests with coverage
	go test ./internal/... -v \
		-count=1  \
		-coverpkg=$(CVPKG) \
		-coverprofile unit-coverage.out \
		--tags=unit

	# Generate human readable coverage result
	go tool cover -html unit-coverage.out -o unit-coverage.html

run-acceptance-tests:
	# Run acceptance tests with coverage
	go test ./internal/application/application_test.go -v \
		-count=1 \
		-coverpkg=$(CVPKG) \
		-coverprofile acceptance-coverage.out \
		--tags=acceptance

	# Generate human readable coverage result
	go tool cover -html acceptance-coverage.out -o acceptance-coverage.html

setup-docker-environment:
	docker-compose up --detach --build

teardown-docker-environment:
	docker-compose down

run-integration-tests:
	# Run integration tests with coverage
	INTEGRATION_REDIS_CONNECTION_STRING=redis://localhost:6379 \
		go test \
		./internal/application/redirection/redis_repository_test.go \
		-v \
		-count=1 \
		-coverpkg=$(CVPKG) \
		-coverprofile integration-coverage-redis.out \
		--tags=integration

	INTEGRATION_MONGO_CONNECTION_STRING=mongodb://root:example@localhost:27017 \
		go test \
		./internal/application/event/mongo_repository_test.go \
		-v \
		-count=1 \
		-coverpkg=$(CVPKG) \
		-coverprofile integration-coverage-mongo.out \
		--tags=integration

	# Generate human readable coverage result
	go tool cover -html integration-coverage-redis.out -o integration-coverage-redis.html
	go tool cover -html integration-coverage-mongo.out -o integration-coverage-mongo.html

run-e2e-tests:
	# Run e2e tests
	E2E_GRPC_SERVER=localhost:6060 \
	go test \
		./internal/port/grpc/server_test.go \
		-v \
		-count=1 \
		--tags=e2e

	E2E_HTTP_SERVER=http://localhost:7070 \
	go test \
		./internal/port/http/server_test.go \
		-v \
		-count=1 \
		--tags=e2e
