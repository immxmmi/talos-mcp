.PHONY: build run tidy test clean docker docker-run

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

docker-run:
	docker run --rm -i \
		$(shell test -f .env && echo --env-file .env) \
		-v $(HOME)/.talos:/root/.talos:ro \
		talos-mcp:latest
