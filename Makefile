PROTO_DIR := internal/pkg/proto
protogen:
	mkdir -p $(PROTO_DIR)/echo
	protoc --proto_path=./proto/echo --go_out=plugins=grpc:$(PROTO_DIR)/echo echo.proto
