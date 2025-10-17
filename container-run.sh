#!/bin/bash
# Local testing script
set -e

IMAGE_NAME="${IMAGE_NAME:-chefly}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
CONTAINER_RUNTIME="${CONTAINER_RUNTIME:-podman}"
CONTAINER_NAME="${CONTAINER_NAME:-chefly}"

if [ -z "$JWT_SECRET" ]; then
    echo "ERROR: JWT_SECRET environment variable is required"
    exit 1
fi

if [ -z "$CLAUDE_API_KEY" ]; then
    echo "ERROR: CLAUDE_API_KEY environment variable is required"
    exit 1
fi

if [ -z "$OPENAI_API_KEY" ]; then
    echo "ERROR: OPENAI_API_KEY environment variable is required"
    exit 1
fi

mkdir -p ./data

${CONTAINER_RUNTIME} run -d \
  --name ${CONTAINER_NAME} \
  --restart unless-stopped \
  --userns=keep-id \
  -p 8080:8080 \
  -v ./data:/app/data:Z \
  -e JWT_SECRET="${JWT_SECRET}" \
  -e DB_PATH="${DB_PATH}" \
  -e CLAUDE_API_KEY="${CLAUDE_API_KEY}" \
  -e CLAUDE_MODEL="${CLAUDE_MODEL:-claude-3-haiku-20240307}" \
  -e OPENAI_API_KEY="${OPENAI_API_KEY}" \
  -e OPENAI_MODEL="${OPENAI_MODEL:-dall-e-3}" \
  -e REGISTRATION_ENABLED="${REGISTRATION_ENABLED:-true}" \
  -e RECIPE_GENERATION_LIMIT="${RECIPE_GENERATION_LIMIT:-unlimited}" \
  -e AUDIT_LOG_ENABLED="${AUDIT_LOG_ENABLED:-true}" \
  -e AUDIT_LOG_LEVEL="${AUDIT_LOG_LEVEL:-info}" \
  -e AUDIT_LOG_FORMAT="${AUDIT_LOG_FORMAT:-json}" \
  --security-opt=no-new-privileges \
  --cap-drop=ALL \
  ${IMAGE_NAME}:${IMAGE_TAG}
