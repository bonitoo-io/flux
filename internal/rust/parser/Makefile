export CC=clang

# Default target
# Does nothing as this used for the Go generate process
all:

clean:
	$(CARGO) clean

build:
	wasm-pack build --dev --scope influxdata

test:
	$(CARGO) test

.PHONY: all clean test
