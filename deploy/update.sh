#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="/opt/dzemo"

if [[ $EUID -eq 0 ]]; then
  echo "Don't run this shell with sudo directly." >&2
  exit 1
fi

echo "Checking latest release..."
LATEST_JSON=$(curl -fsSL "https://api.github.com/repos/zijoziks/dzemo/releases/latest")
LATEST_TAG=$(echo "$LATEST_JSON" | jq -r '.tag_name')
DOWNLOAD_URL=$(echo "$LATEST_JSON" | jq -r '.assets[] | select(.name == "Dzemo") | .browser_download_url')

CURRENT_TAG=""
[[ -f "$INSTALL_DIR/VERSION" ]] && CURRENT_TAG=$(cat "$INSTALL_DIR/VERSION")

if [[ "$LATEST_TAG" == "$CURRENT_TAG" ]]; then
  echo "Already up to date ($CURRENT_TAG)."
  exit 0
fi

echo "Updating $CURRENT_TAG -> $LATEST_TAG"
curl -fsSL -o /tmp/Dzemo "$DOWNLOAD_URL"
chmod +x /tmp/Dzemo

sudo mv /tmp/Dzemo "$INSTALL_DIR/Dzemo"
sudo chown dzemo:dzemo "$INSTALL_DIR/Dzemo"
echo "$LATEST_TAG" | sudo tee "$INSTALL_DIR/VERSION" > /dev/null
sudo chown dzemo:dzemo "$INSTALL_DIR/VERSION"

sudo systemctl restart dzemo
sudo systemctl status dzemo --no-pager
