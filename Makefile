.PHONY: all setup sysconf

all: setup sysconf homeconf

setup: bin/setup

sysconf: bin/sysconf

homeconf: bin/homeconf

bin/setup: setup/setup.go
	cd setup && go build -ldflags="-s -w" -o ../bin/setup

bin/sysconf: sysconf/sysconf.go
	cd sysconf && go build -ldflags="-s -w" -o ../bin/sysconf

bin/homeconf: homeconf/homeconf.go
	cd homeconf && go build -ldflags="-s -w" -o ../bin/homeconf
