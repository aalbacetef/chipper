all: build-emu build-dumprom 


build: web 

lint: web-lint go-lint

test: web-test go-test

### Go-only build tasks 

go-fmt:
	goimports -w .

go-test: go-fmt
	go test -v ./...

go-lint: go-fmt
	golangci-lint run 

mk-bin-dir:
	mkdir -p ./bin/ 

build-emu: go-fmt mk-bin-dir
	go build -o bin/ ./cmd/emu/ 

build-dumprom: go-fmt mk-bin-dir
	go build -o bin/ ./cmd/dumprom/


## WebUI tasks 


### linting and testing 

web-test:
	cd ./webui && bun x vitest

type-check:
	cd webui && bun x vue-tsc --build --force

web-prettier:
	cd webui && bun x prettier --write ./src

web-lint: type-check web-prettier

### bundling 

make-manifest:
	./roms/manifest.fish

copy-roms: make-manifest 
	rm -rf ./webui/public/roms/
	mkdir ./webui/public/roms
	cp -r ./roms/set-* ./webui/public/roms/
	cp ./roms/manifest.json ./webui/public/roms/


build-wasm: export GOOS=js
build-wasm: export GOARCH=wasm 
build-wasm: mk-bin-dir fmt
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/webui.wasm ./cmd/webui/ 

copy-wasm: build-wasm
	mkdir -p ./webui/public
	cp ./bin/webui.wasm ./webui/public/


### WebUI build

clean-web: 
	rm -rf ./webui/dist/

web: copy-wasm copy-roms
	cd webui && bun x vite build

web-local: clean-web copy-wasm copy-roms
	cd webui && bun x vite build -m dev 

web-serve: web-local
	http-server ./webui/dist/

dev: fmt
	cd webui && bun x vite 


.PHONY: build build-emu build-dumprom lint dev test mk-bin-dir fmt 
.PHONY: web copy-wasm copy-roms make-manifest web-test

