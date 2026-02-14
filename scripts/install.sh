#!/bin/bash
set -euo pipefail

REPO="crypto-cube/dsus"
BIN_PATH="/usr/bin/dsus"
CONF_DIR="/etc/dsus"
SERVICE_FILE="/etc/systemd/system/dsus.service"

if [ "$(id -u)" -ne 0 ]; then
    echo "This script must be run as root"
    exit 1
fi

# Install prerequisites for easy-wg-quick
echo ""
echo "Installing WireGuard and prerequisites..."
apt-get update -qq
apt-get install -y -qq wireguard wireguard-tools qrencode

# Prompt for public key
echo ""
read -rp "Path to public key (PEM format): " PUBKEY_PATH < /dev/tty
if [ ! -f "$PUBKEY_PATH" ]; then
    echo "File not found: ${PUBKEY_PATH}, aborting."
    exit 1
fi
PUBKEY=$(cat "$PUBKEY_PATH")

# Prompt for basic auth credentials
echo ""
read -rp "Basic auth username (leave empty to disable): " AUTH_USER < /dev/tty
if [ -n "$AUTH_USER" ]; then
    read -rsp "Basic auth password: " AUTH_PASS < /dev/tty
    echo ""
    if [ -z "$AUTH_PASS" ]; then
        echo "Password cannot be empty when username is set, aborting."
        exit 1
    fi
fi

# Download latest release
echo ""
echo "Downloading latest release..."
DOWNLOAD_URL=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep -o "https://github.com/${REPO}/releases/download/[^\"]*" \
    | head -1)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Could not find latest release, aborting."
    exit 1
fi

curl -fsSL -o "$BIN_PATH" "$DOWNLOAD_URL"
chmod 755 "$BIN_PATH"

# Create system user
if ! id "dsus" &>/dev/null; then
    useradd --system --no-create-home --shell /usr/sbin/nologin dsus
fi

# Set up directories
mkdir -p "${CONF_DIR}/certs"
mkdir -p /var/lib/dsus/files
mkdir -p /var/lib/dsus/wg
echo "$PUBKEY" > "${CONF_DIR}/certs/publickey.pub"

# Download easy-wg-quick
curl -fsSL -o /var/lib/dsus/wg/easy-wg-quick \
    "https://raw.githubusercontent.com/burghardt/easy-wg-quick/master/easy-wg-quick"
chmod 755 /var/lib/dsus/wg/easy-wg-quick

chown -R dsus:dsus "$CONF_DIR"
chown -R dsus:dsus /var/lib/dsus

# Initialize easy-wg-quick
echo ""
echo "Initializing WireGuard via easy-wg-quick..."
cd /var/lib/dsus/wg
./easy-wg-quick

# Prompt for devices prefix
echo ""
read -rp "Devices prefix: " DEVICES_PREFIX < /dev/tty
if [ -z "$DEVICES_PREFIX" ]; then
    echo "No devices prefix provided, aborting."
    exit 1
fi

# Build environment file
ENV_FILE="${CONF_DIR}/env"
: > "$ENV_FILE"
if [ -n "${AUTH_USER:-}" ]; then
    echo "DSUS_USER=${AUTH_USER}" >> "$ENV_FILE"
    echo "DSUS_PASS=${AUTH_PASS}" >> "$ENV_FILE"
fi
echo "DSUS_DEVICES_PREFIX=${DEVICES_PREFIX}" >> "$ENV_FILE"
chmod 600 "$ENV_FILE"
chown dsus:dsus "$ENV_FILE"

# Install systemd service
cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=Darn Simple Update Server

[Service]
User=dsus
Group=dsus
Restart=on-failure
EnvironmentFile=${CONF_DIR}/env
ExecStart=${BIN_PATH}

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now dsus

echo ""
echo "dsus installed and running."
echo "Check status: systemctl status dsus"
