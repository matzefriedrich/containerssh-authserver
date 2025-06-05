#!/bin/sh
set -e

# Defaults (override with env vars or CLI args if needed)
CI_COMMIT_SHORT_SHA="${CI_COMMIT_SHORT_SHA:-$(git rev-parse --short=8 HEAD)}"
APP_RELEASE_DATE="$(date -u +%Y-%m-%d)"
APP_RELEASE="${APP_RELEASE:-matzefriedrich/containerssh-authserver}"
APP_VERSION="${APP_VERSION:-0.1.0}"
IMAGE_NAME="${IMAGE_NAME:-mobymatze/containerssh-authserver}"

echo "Building Docker image with:"
echo "  Commit hash:  $CI_COMMIT_SHORT_SHA"
echo "  Release date: $APP_RELEASE_DATE"
echo "  Release:      $APP_RELEASE"
echo "  Version:      $APP_VERSION"
echo "  Image name:   $IMAGE_NAME"

docker build \
    --build-arg CI_COMMIT_SHORT_SHA="$CI_COMMIT_SHORT_SHA" \
    --build-arg APP_RELEASE_DATE="$APP_RELEASE_DATE" \
    --build-arg APP_RELEASE="$APP_RELEASE" \
    --build-arg APP_VERSION="$APP_VERSION" \
  -t "$IMAGE_NAME:$APP_VERSION" \
  -f Dockerfile .
