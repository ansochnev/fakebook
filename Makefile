.PHONY: all clean fakebook

all: fakebook

clean:
	rm -f fakebook

fakebook:
	go build -o fakebook -tags netgo cmd/fakebook/*

tests:
	go test ./...
