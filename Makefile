.PHONY: setup

all: setup

setup: bin/setup

bin/setup: setup/setup.go
	cd setup && go build -ldflags="-s -w" -o ../bin/setup
