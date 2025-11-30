.PHONY: build test run clean

APP_NAME=kafctl
VERSION=0.0.1
PROJECT=kafctl
WORK_DIR=.

build:
	@echo "Building $(APP_NAME) version $(VERSION)..."
	go build -o ${APP_NAME} ${WORK_DIR}/cmd/kafctl

build-windows:
	@echo "Building $(APP_NAME) for Windows..."
	docker build -t kafctl-windows -f Dockerfile.windows .
	docker run --rm -v "$$(pwd):/app" kafctl-windows
	@echo "Windows build complete. Check build.log for details."

test:
	go test ${WORK_DIR}/...

run: build
	./kafctl

clean:
	@echo "Deleting $(APP_NAME) version $(VERSION)..."
	rm -rf ${APP_NAME}