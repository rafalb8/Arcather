LDFLAGS = "-s -w"

.PHONY: all linux windows

all: linux windows

linux:
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=linux \
	go build -ldflags=$(LDFLAGS) -o bin/arcather ./cmd/arcather
	strip bin/arcather

windows:
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=windows \
	go build -ldflags=$(LDFLAGS) -o bin/arcather.exe ./cmd/arcather
	strip bin/arcather.exe

clean:
	rm -rf bin/*
