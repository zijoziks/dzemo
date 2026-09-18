#!/usr/bin/env bash
set -euo pipefail

REPO="zijoziks/dzemo"
INSTALL_DIR="/opt/dzemo"
DEPLOY_DIR="$INSTALL_DIR/deploy"

if [[ $EUID -eq 0 ]]; then
  echo "Please don't run this script as root." >&2
  exit 1
fi

sudo -v

echo "Fetching latest release..."
LATEST_JSON=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest")
LATEST_TAG=$(echo "$LATEST_JSON" | jq -r '.tag_name')
DOWNLOAD_URL=$(echo "$LATEST_JSON" | jq -r '.assets[] | select(.name == "Dzemo") | .browser_download_url')

echo "Installing $LATEST_TAG"
curl -fsSL -o /tmp/Dzemo "$DOWNLOAD_URL"
chmod +x /tmp/Dzemo

RAW_BASE="https://raw.githubusercontent.com/$REPO/$LATEST_TAG"
curl -fsSL -o /tmp/dzemo.service "$RAW_BASE/deploy/dzemo.service"
curl -fsSL -o /tmp/update.sh "$RAW_BASE/deploy/update.sh"
curl -fsSL -o /tmp/uninstall.sh "$RAW_BASE/deploy/uninstall.sh"
curl -fsSL -o /tmp/README.md "$RAW_BASE/README.md"
chmod +x /tmp/update.sh
chmod +x /tmp/uninstall.sh

id -u dzemo &>/dev/null || sudo useradd --system --no-create-home --shell /usr/sbin/nologin dzemo

sudo mkdir -p "$INSTALL_DIR" "$DEPLOY_DIR"
sudo mv /tmp/Dzemo "$INSTALL_DIR/Dzemo"
sudo chown dzemo:dzemo "$INSTALL_DIR/Dzemo"
echo "$LATEST_TAG" | sudo tee "$INSTALL_DIR/VERSION" > /dev/null
sudo chown dzemo:dzemo "$INSTALL_DIR/VERSION"

sudo mv /tmp/dzemo.service /etc/systemd/system/dzemo.service
sudo mv /tmp/update.sh "$DEPLOY_DIR/update.sh"
sudo mv /tmp/uninstall.sh "$DEPLOY_DIR/uninstall.sh"
sudo mv /tmp/README.md "$DEPLOY_DIR/README.md"

sudo mkdir -p /etc/dzemo
read -rs -p "Discord bot token: " TOKEN < /dev/tty
echo
printf '%s' "$TOKEN" | sudo systemd-creds encrypt --name=dzemo_token - /etc/dzemo/dzemo_token.cred
unset TOKEN

read -rp "Guild (server) ID: " GUILD_ID < /dev/tty
read -rp "Welcome channel ID: " WELCOME_CHANNEL_ID < /dev/tty
read -rp "Member role ID (granted on /verify): " MEMBER_ROLE_ID < /dev/tty

sudo mkdir -p /var/lib/dzemo
cat <<EOF | sudo tee /var/lib/dzemo/config.json > /dev/null
{
    "guildID": "$GUILD_ID",
    "welcomeChannelID": "$WELCOME_CHANNEL_ID",
    "memberRoleID": "$MEMBER_ROLE_ID"
}
EOF
sudo chown -R dzemo:dzemo /var/lib/dzemo
sudo chmod 600 /var/lib/dzemo/config.json

sudo systemctl daemon-reload
sudo systemctl enable --now dzemo
sudo systemctl status dzemo --no-pager

echo
echo "Installed $LATEST_TAG. Deployment files for future reference are in $DEPLOY_DIR (update.sh, README.md, uninstall.sh)."