#!/usr/bin/env bash
set -euo pipefail

if [[ $EUID -eq 0 ]]; then
  echo "Please don't run this script as root." >&2
  exit 1
fi

read -rp "This will stop dzemo and remove the service, its user and everything else tied to it. Continue? [y/N] " CONFIRM < /dev/tty
if [[ "$CONFIRM" != "y" && "$CONFIRM" != "Y" ]]; then
  echo "Aborted."
  exit 0
fi

sudo -v

echo "Stopping and disabling service..."
sudo systemctl disable --now dzemo || true

echo "Removing unit file..."
sudo rm -f /etc/systemd/system/dzemo.service
sudo systemctl daemon-reload

echo "Removing token credential..."
sudo rm -rf /etc/dzemo

echo "Removing dzemo user..."
if id -u dzemo &>/dev/null; then
  sudo userdel dzemo
fi

echo "Removing /var/lib/dzemo..."
sudo rm -rf /var/lib/dzemo

echo "Removing /opt/dzemo..."
sudo rm -rf /opt/dzemo

echo "Done. Dzemo is gone now!"