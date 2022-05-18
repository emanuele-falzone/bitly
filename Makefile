download:
	curl --create-dirs  -o google/api/annotations.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto
	curl --create-dirs  -o google/api/http.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto

proto:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		 	google/api/annotations.proto \
         	google/api/http.proto \
			internal/adapter/service/grpc/pb/bitly_service.proto

gateway:
	protoc  --grpc-gateway_out=. \
			--grpc-gateway_opt=paths=source_relative \
           	internal/adapter/service/grpc/pb/bitly_service.proto

openapidocs: 
	protoc \
    --openapiv2_out . \
    --openapiv2_opt logtostderr=true \
    internal/adapter/service/grpc/pb/bitly_service.proto

protodocs:
	protoc \
	--doc_out=./docs \
	internal/adapter/service/grpc/pb/bitly_service.proto
