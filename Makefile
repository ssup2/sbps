SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 0.1.0
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean install uninstall fmt run

all: install

$(TARGET): $(SRC)
	@go build $(LDFLAGS) -o $(TARGET) $(SRC)

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET)

install:
	@go install $(LDFLAGS) ./cmd/$(TARGET)

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

run: install
	@$(TARGET)
