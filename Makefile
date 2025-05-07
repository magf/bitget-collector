.PHONY: all build run test clean deb

all: build

# Build the collector binary
build:
	go build -o collector ./cmd/collector

# Run the collector with a default pair for testing
run: build
	./collector -pair=BTCUSDT -debug

# Run tests (placeholder, add tests to cmd/collector if needed)
test:
	go test ./...

# Clean up generated files
clean:
	rm -f collector bitget-collector.deb
	rm -rf deb-package

# Build the DEB package
deb: build
	./scripts/build-deb.sh
