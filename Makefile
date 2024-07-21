BINARY_NAME={project-name}
# GOPRIVATE=gitlab.revenuemonster.my/dinar-wallet/*

all: test build
test:
	go test -v ./...
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags timetzdata -ldflags="-w -s" -o $(BINARY_NAME) .
generate-protobuf:
	@rm -rf protobuf/go/** && protoc -I=protobuf/proto --go_out=paths=source_relative:protobuf/go --go-grpc_out=paths=source_relative:protobuf/go protobuf/proto/**/*.proto