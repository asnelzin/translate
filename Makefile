SHELL := /bin/bash

%:
	@:

binary-osx:
	./scripts/build/osx.sh

build: binary-osx

clean:
	rm -rf ./target/*

tests:
	go test -v github.com/asnelzin/translate/yandex

coverage:
	go test -cover -coverprofile=coverage.txt github.com/asnelzin/translate/yandex

coverage-html:
	go tool cover -html=coverage.txt

.PHONY: tests