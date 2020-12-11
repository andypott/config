.PHONY: all setup sysconf

all: setup sysconf

setup: bin/setup

sysconf: bin/sysconf

bin/setup: setup/setup.go
	cd setup && go build -ldflags="-s -w" -o ../bin/setup

bin/sysconf: sysconf/sysconf.go
	cd sysconf && go build -ldflags="-s -w" -o ../bin/sysconf
