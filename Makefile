.PHONY: build run tidy test clean docker

build:
	cd src && go build -o ../bin/talos-mcp .

run: build
	./bin/talos-mcp

tidy:
	cd src && go mod tidy

test:
	cd src && go test ./...

clean:
	rm -rf bin/

docker:
	docker build -t talos-mcp:latest .
