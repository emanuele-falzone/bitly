generate-repositoty-mock:
	go install github.com/golang/mock/mockgen
	mockgen -destination=./test/mock/redirection_repository.go \
		-package=mock \
		-mock_names=Repository=MockRedirectionRepository \
		github.com/emanuelefalzone/bitly/internal/domain/redirection Repository

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
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
	protoc \
		--doc_out=./docs \
		internal/adapter/service/grpc/pb/bitly_service.proto