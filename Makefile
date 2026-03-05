BINARY_NAME=outline
VERSION=0.1.0

.PHONY: build clean install run

build:
	go build -ldflags="-X outline-cli/cmd.Version=$(VERSION)" -o $(BINARY_NAME) ./main.go

clean:
	rm -f $(BINARY_NAME)

install: build
	cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

run: build
	./$(BINARY_NAME)

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy
