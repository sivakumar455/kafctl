.PHONY: build test run clean

APP_NAME=kafctl
VERSION=0.0.1
PROJECT=kafctl
WORK_DIR=.

build:
	@echo "Building $(APP_NAME) version $(VERSION)..."
	go build -o ${APP_NAME} ${WORK_DIR}/cmd/kafctl

test:
	go test ${WORK_DIR}/...

run: build
	./kafctl

clean:
	@echo "Deleting $(APP_NAME) version $(VERSION)..."
	rm -rf ${APP_NAME}