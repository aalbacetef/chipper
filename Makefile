all: build-emu build-dumprom 

fmt:
	goimports -w .

mk-bin-dir:
	mkdir -p ./bin/ 

build-emu: fmt mk-bin-dir
	go build -o bin/ ./cmd/emu/ 

build-dumprom: fmt mk-bin-dir
	go build -o bin/ ./cmd/dumprom/

test: fmt
	go test -v ./...

lint: fmt
	golangci-lint run 
