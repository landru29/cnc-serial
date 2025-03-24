.PHONY: cnc-serial cnc-serial-minimal all build-debug-serial install-tools debug kill-debug

all: cnc-serial cnc-serial-minimal

cnc-serial:
	go build --tags withbutton -o cnc-serial ./cmd/...

cnc-serial-minimal:
	go build -o cnc-serial-minimal ./cmd/...
clean:
	rm -f cnc-serial cnc-serial-minimal __debug_bin

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

