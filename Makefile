generate-repository-mock:
	go install github.com/golang/mock/mockgen
	mockgen -destination=./test/mock/redirection_repository.go \
		-package=mock \
		-mock_names=Repository=MockRedirectionRepository \
		github.com/emanuelefalzone/bitly/internal/domain/redirection Repository

generate-key-generator-mock:
	go install github.com/golang/mock/mockgen
	mockgen -destination=./test/mock/key_generator_service.go \
		-package=mock \
		github.com/emanuelefalzone/bitly/internal/service KeyGenerator

generate-grpc-server:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
			internal/adapter/service/grpc/pb/bitly_service.proto

generate-grpc-server-documentation:
	mkdir -p ./docs
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
	protoc \
		--doc_out=./docs \
		internal/adapter/service/grpc/pb/bitly_service.proto

build: generate-grpc-server
	go build -v ./cmd/main.go

run-unit-tests: generate-repository-mock generate-key-generator-mock
	CVPKG=$(go list ./internal/... | grep -v pb | tr '\n' ',') 
	go test ./internal/... -coverpkg=$CVPKG -coverprofile coverage.out -v 
	go tool cover -html coverage.out -o coverage.html

run-integration-tests:
	docker-compose -f test/integration/docker-compose.yml up -d
	sleep 5
	INTEGRATION_REDIS_CONNECTION_STRING=redis://localhost:6379 \
		go test ./test/integration/integration_test.go -v
	docker-compose -f test/integration/docker-compose.yml down

run-acceptance-tests:
	docker-compose up -d
	sleep 5
	ACCEPTANCE_REDIS_CONNECTION_STRING=redis://localhost:6379 \
		ACCEPTANCE_GRPC_SERVER=localhost:4000 \
		go test ./test/acceptance/go_acceptance_test.go -v
	
	ACCEPTANCE_REDIS_CONNECTION_STRING=redis://localhost:6379 \
		ACCEPTANCE_GRPC_SERVER=localhost:4000 \
		go test ./test/acceptance/grpc_acceptance_test.go -v

	docker-compose down