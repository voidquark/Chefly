#!/bin/bash
# Local Testing Script
set -e

IMAGE_NAME="${IMAGE_NAME:-chefly}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
CONTAINER_RUNTIME="${CONTAINER_RUNTIME:-podman}"

USER_UID="${USER_UID:-1000}"
USER_GID="${USER_GID:-1000}"

echo "🏗️  Building Chefly container image"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📦 Image: ${IMAGE_NAME}:${IMAGE_TAG}"
echo "👤 User: UID=${USER_UID} GID=${USER_GID}"
echo "🔧 Runtime: ${CONTAINER_RUNTIME}"
echo ""

# Build container
${CONTAINER_RUNTIME} build \
  --build-arg USER_UID=${USER_UID} \
  --build-arg USER_GID=${USER_GID} \
  --tag ${IMAGE_NAME}:${IMAGE_TAG} \
  -f Dockerfile \
  .

echo ""
echo "✅ Build complete!"
echo ""
echo "📊 Image details:"
${CONTAINER_RUNTIME} images ${IMAGE_NAME}:${IMAGE_TAG}
