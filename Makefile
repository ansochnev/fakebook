.PHONY: all clean fakebook

all: fakebook

clean:
	rm -f fakebook

fakebook:
	go build -o fakebook cmd/fakebook/main.go
