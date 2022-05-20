generate-code:
	# Install mockgen
	go install github.com/golang/mock/mockgen

	# Generate mocks for test purposes
	# Generate mock redirection repository
	mockgen -destination=./test/mock/redirection_repository.go \
		-package=mock \
		-mock_names=Repository=MockRedirectionRepository \
		github.com/emanuelefalzone/bitly/internal/domain/redirection Repository

	# Generate mock event repository
	mockgen -destination=./test/mock/event_repository.go \
		-package=mock \
		-mock_names=Repository=MockEventRepository \
		github.com/emanuelefalzone/bitly/internal/domain/event Repository

	# Generate mock key generator
	mockgen -destination=./test/mock/key_generator_service.go \
		-package=mock \
		github.com/emanuelefalzone/bitly/internal/service KeyGenerator

	# Generate mock event listener
	mockgen -destination=./test/mock/event_listener.go \
		-package=mock \
		github.com/emanuelefalzone/bitly/internal/domain/event Listener

	# Install protobuf and grpc
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/protobuf/cmd/protoc-gen-go

	# Generate protobuf and grpc server
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		internal/adapter/service/grpc/pb/bitly_service.proto

generate-docs:
	# Ensure docs directory exists
	mkdir -p ./docs

	# Install proto documentation generator
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

	# Generate documentation for grpc service
	protoc \
		--doc_out=./docs \
		internal/adapter/service/grpc/pb/bitly_service.proto

build: generate-code
	go build -v ./cmd/main.go

run-unit-tests: generate-code
	CVPKG=$(go list ./internal/... | grep -v pb | tr '\n' ',')
	go test ./internal/... -coverpkg=$(CVPKG) -coverprofile coverage.out -v --tags=unit
	go tool cover -html coverage.out -o coverage.html

run-acceptance-tests: generate-code
	go test ./internal/... -v --tags=acceptance

run-integration-tests: generate-code
	docker-compose up -d --build --force-recreate
	sleep 5
	INTEGRATION_REDIS_CONNECTION_STRING=redis://localhost:6379 \
	INTEGRATION_MONGO_CONNECTION_STRING=mongodb://root:example@localhost:27017 \
		go test ./internal/... --tags=integration -v
	docker-compose down



run-e2e-tests: generate-code
	docker-compose up -d --build --force-recreate
	sleep 5
	E2E_GRPC_SERVER=localhost:4000 \
		go test ./internal/... -v --tags=acceptance
	docker-compose down