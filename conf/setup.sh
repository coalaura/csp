#!/bin/bash

set -euo pipefail

# --- HARDWARE PERMISSIONS ---
# If this service requires hardware access, you likely need a udev rule
# to assign ownership to the 'csp' user.
# Example: /etc/udev/rules.d/99-csp.rules
# SUBSYSTEM=="usb", ATTRS{idVendor}=="XXXX", OWNER="csp"
# ----------------------------

echo "Linking sysusers config..."

mkdir -p /etc/sysusers.d

if [ -f /etc/sysusers.d/csp.conf ]; then
    rm /etc/sysusers.d/csp.conf
fi

ln -s "/var/example.com/conf/csp.conf" /etc/sysusers.d/csp.conf

echo "Creating user..."

systemd-sysusers

echo "Linking unit..."

if [ -f /etc/systemd/system/csp.service ]; then
    rm /etc/systemd/system/csp.service
fi

systemctl link "/var/example.com/conf/csp.service"

if command -v logrotate >/dev/null 2>&1; then
    echo "Linking logrotate config..."

    if [ -f /etc/logrotate.d/csp ]; then
        rm /etc/logrotate.d/csp
    fi

    ln -s "/var/example.com/conf/csp_logs.conf" /etc/logrotate.d/csp
else
    echo "Logrotate not found, skipping..."
fi

echo "Reloading daemon..."

systemctl daemon-reload
systemctl enable csp

echo "Fixing initial permissions..."

mkdir -p "/var/example.com/logs"

chown -R csp:csp "/var/example.com"

find "/var/example.com" -type d -exec chmod 755 {} +
find "/var/example.com" -type f -exec chmod 644 {} +

chmod +x "/var/example.com/csp"

echo "Setup complete, starting service..."

service csp restart

echo "Done."
