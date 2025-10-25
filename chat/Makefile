LOCAL_BIN:=$(CURDIR)/bin

# установка локально
install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest

# установка глобально
get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-chat-api

# generate-chat-api:
# 	mkdir -p pkg/api/chat_v1
# 	protoc --proto_path api/chat_v1 \
# 	--go_out=pkg/api/chat_v1 --go_opt=paths=source_relative \
# 	--plugin=protoc-gen-go=bin/protoc-gen-go \
# 	--go-grpc_out=pkg/api/chat_v1 --go-grpc_opt=paths=source_relative \
# 	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
# 	api/chat_v1/chat.proto

generate-chat-api:
	mkdir -p pkg/api/chat_v1
	protoc --proto_path api/chat_v1 --proto_path vendor.protogen \
	--go_out=pkg/api/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/api/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/api/chat_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	api/chat_v1/chat.proto


vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi