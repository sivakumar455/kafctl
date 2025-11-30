#!/bin/bash

# Image name
IMAGE_NAME="kafctl-windows"

echo "Building Docker image..."
docker build -t $IMAGE_NAME -f Dockerfile.windows .

if [ $? -ne 0 ]; then
    echo "Docker image build failed."
    exit 1
fi

echo "Running build in container..."
# Mount current directory to /app
docker run --rm -v "$(pwd):/app" $IMAGE_NAME

if [ $? -eq 0 ]; then
    echo "Build completed successfully. Check build.log for details."
else
    echo "Build failed. Check build.log for details."
    exit 1
fi
