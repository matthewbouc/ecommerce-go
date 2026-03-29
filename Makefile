server:
	nodemon --watch './**/*go' --signal SIGTERM --exec APP_ENV=dev 'go' run main.go

run:
	APP_ENV=dev 'go' run main.go

lint:
	golangci-lint run ./...

# Generate Go code from all .proto files.
# --go_out / --go-grpc_out: where to write the files (. = project root)
# paths=source_relative: output mirrors the proto file's directory structure
# Run this any time a .proto file changes.
proto:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		proto/catalog/catalog.proto