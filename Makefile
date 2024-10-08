all: build-emu build-dumprom 

fmt:
	goimports -w .

mk-bin-dir:
	mkdir -p ./bin/ 

build-emu: fmt mk-bin-dir
	go build -o bin/ ./cmd/emu/ 

build-dumprom: fmt mk-bin-dir
	go build -o bin/ ./cmd/dumprom/

 
build-wasm: export GOOS=js
build-wasm: export GOARCH=wasm 
build-wasm: mk-bin-dir fmt
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/webui.wasm ./cmd/webui/ 

copy-wasm: build-wasm
	cp ./bin/webui.wasm ./webui/public/

web: copy-wasm
	cd webui && bun x vite build

test: fmt test-go
test-go: fmt
	go test -v ./...

lint: fmt
	golangci-lint run 

dev: fmt
	cd webui && bun x vite 


type-check:
	cd webui && bun x vue-tsc --build --force


.PHONY: build-emu build-dumprom lint dev test mk-bin-dir fmt web copy-wasm 

