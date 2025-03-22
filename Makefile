.PHONY: cnc-serial cnc-serial-minimal all build-debug-serial

all: cnc-serial cnc-serial-minimal

cnc-serial:
	go build --tags withbutton -o cnc-serial ./cmd/...

cnc-serial-minimal:
	go build -o cnc-serial-minimal ./cmd/...
clean:
	rm -f cnc-serial cnc-serial-minimal __debug_bin

lint:
	golangci-lint run ./...

debug: __debug_bin
	dlv --listen=:2345 --headless=true --api-version=2 --log exec __debug_bin

__debug_bin:
	go build -gcflags="-N -l" -a -o __debug_bin ./cmd