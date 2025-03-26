.PHONY: cnc cnc-minimal all build-debug install-tools debug kill-debug __debug_bin

all: cnc cnc-minimal

cnc:
	go build --tags withbutton -o cnc ./cmd/...

cnc-minimal:
	go build -o cnc-minimal ./cmd/...
clean:
	rm -f cnc cnc-minimal __debug_bin

lint:
	golangci-lint run ./...

install-tools: go install github.com/go-delve/delve/cmd/dlv@latest

debug: __debug_bin
	dlv \
	  --listen=:2345 \
	  --log=true \
	  --headless=true \
	  --accept-multiclient \
	  --api-version=2 \
	  exec __debug_bin -- --dry-run internal/gcode/grbl/testdata/prog01.gcode

kill-debug:
	kill `ps aux | grep "dlv" | grep __debug_bin | awk '{print $$2}'`

__debug_bin:
	go build -gcflags "all=-N -l" -o __debug_bin ./cmd

