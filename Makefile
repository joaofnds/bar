GIT_COMMIT=$(shell git rev-parse --short HEAD)

.PHONY: test

test:
	go test

build:
	docker build -t joaofnds/bar:$(GIT_COMMIT) .

push:
	docker push joaofnds/bar:$(GIT_COMMIT)