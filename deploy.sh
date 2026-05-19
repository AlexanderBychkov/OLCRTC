#!/bin/bash
set -e

if [ ! -f "./.env" ]; then
  echo "ERROR: .env not found. Copy .env.example to .env and fill in values."
  exit 1
fi
source "./.env"

IMAGE=olcrtc/server:local

deploy_user() {
  local N=
  local CLIENT_ID=
  local KEY=

  echo "==> Deploying olcrtc-wb-"
  docker stop "olcrtc-wb-" 2>/dev/null || true
  docker rm   "olcrtc-wb-" 2>/dev/null || true
  docker run -d --name "olcrtc-wb-" --restart unless-stopped --init     -e OLCRTC_CARRIER=wbstream     -e OLCRTC_TRANSPORT=vp8channel     -e OLCRTC_VP8_FPS=60     -e OLCRTC_VP8_BATCH=64     -e OLCRTC_CLIENT_ID=""     -e OLCRTC_KEY=""     -e OLCRTC_DNS=1.1.1.1:53     -e OLCRTC_WB_USER_TOKEN=""     -v "olcrtc-state-wb-:/var/lib/olcrtc"     ""
  echo "    started"
}

echo "==> Building "
docker build -t "" "."

deploy_user 1 "" ""
deploy_user 2 "" ""

echo "==> Done"
sleep 5
docker logs olcrtc-wb-1 2>&1 | tail -3
docker logs olcrtc-wb-2 2>&1 | tail -3
