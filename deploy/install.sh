#!/usr/bin/env bash
set -euo pipefail

if [[ $EUID -eq 0 ]]; then
  echo "Please don't run the script as root." >&2
  exit 1
fi

if [[ ! -f deploy/dzemo.service ]]; then
  echo "Run this from the project root please!"
  exit 1
fi

sudo -v

echo "Fetching latest release..."
LATEST_JSON=$(curl -fsSL "https://api.github.com/repos/zijoziks/dzemo/releases/latest")
LATEST_TAG=$(echo "$LATEST_JSON" | jq -r '.tag_name')
DOWNLOAD_URL=$(echo "$LATEST_JSON" | jq -r '.assets[] | select(.name == "Dzemo") | .browser_download_url')

echo "Installing $LATEST_TAG"
curl -fsSL -o /tmp/Dzemo "$DOWNLOAD_URL"
chmod +x /tmp/Dzemo

id -u dzemo &>/dev/null || sudo useradd --system --no-create-home --shell /usr/sbin/nologin dzemo

sudo mkdir -p /opt/dzemo
sudo mv /tmp/Dzemo /opt/dzemo/Dzemo
sudo chown dzemo:dzemo /opt/dzemo/Dzemo
echo "$LATEST_TAG" | sudo tee /opt/dzemo/VERSION > /dev/null
sudo chown dzemo:dzemo /opt/dzemo/VERSION

sudo cp deploy/dzemo.service /etc/systemd/system/dzemo.service

sudo mkdir -p /etc/dzemo
read -rs -p "Discord bot token: " TOKEN
echo
printf '%s' "$TOKEN" | sudo systemd-creds encrypt --name=dzemo_token - /etc/dzemo/dzemo_token.cred
unset TOKEN

read -rp "Guild (server) ID: " GUILD_ID
read -rp "Welcome channel ID: " WELCOME_CHANNEL_ID
read -rp "Member role ID (granted on /verify): " MEMBER_ROLE_ID

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