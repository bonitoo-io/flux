export CC=clang

# Default target, build
all:
	wasm-pack build --dev --scope influxdata

clean:
	$(CARGO) clean

test:
	$(CARGO) test

.PHONY: all clean test
