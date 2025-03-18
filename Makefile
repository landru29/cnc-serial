cnc-serial:
	go build --tags withbutton -o cnc-serial ./cmd/...

cnc-serial-minimal:
	go build -o cnc-serial-minimal ./cmd/...
clean:
	rm -f cnc-serial cnc-serial-minimal

lint:
	golangci-lint run ./...