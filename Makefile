BINARY_PATH=./bin/acr
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
MAIN_FOLDER=./cmd/acr/

clean:
	rm -f $(BINARY_PATH)

build:
	$(GOBUILD) -o $(BINARY_PATH) $(MAIN_FOLDER)

test:
	$(GOTEST) -v -vet=off $(MAIN_FOLDER)

all: clean test build

.DEFAULT_GOAL := all
.PHONY: all clean build test
