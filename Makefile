all: build-emu build-dumprom 

mk-bin-dir:
	mkdir -p ./bin/ 

build-emu: mk-bin-dir
	go build -o bin/ ./cmd/emu/ 

build-dumprom: mk-bin-dir
	go build -o bin/ ./cmd/dumprom/

test: 
	go test -v ./...

lint:
	golangci-lint run 
