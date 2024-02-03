.PHONY: all generate buildrom run

all: generate buildrom run

run:
	go run .

generate:
	go build -o ./build/protoc-gen-go-wasm ./cmd/protoc-gen-go-wasm
	rm -rf internal/gen
	buf generate

buildrom:
	rm -f ./examples/basic-rom/main.wasm
	GOOS=wasip1 GOARCH=wasm go build -o ./examples/basic-rom/main.wasm ./examples/basic-rom/main.go

buildhw:
	rm -f ./examples/hello-world/main.wasm
	GOOS=wasip1 GOARCH=wasm go build -o ./examples/hello-world/main.wasm ./examples/hello-world/main.go
