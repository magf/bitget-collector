#!/bin/bash

# Build the Go binary
go build -o collector ./cmd/collector

# Create temporary directory for DEB package
mkdir -p deb-package/usr/bin
mkdir -p deb-package/etc/systemd/system
mkdir -p deb-package/var/lib/bitget-collector
mkdir -p deb-package/DEBIAN

# Copy files to DEB structure
cp collector deb-package/usr/bin/
cp debian/control deb-package/DEBIAN/
cp debian/preinst deb-package/DEBIAN/
cp debian/postinst deb-package/DEBIAN/
cp debian/collector@.service deb-package/etc/systemd/system/

# Set permissions
chmod 755 deb-package/usr/bin/collector
chmod 644 deb-package/etc/systemd/system/collector@.service
chmod 755 deb-package/DEBIAN/preinst
chmod 755 deb-package/DEBIAN/postinst
chmod 644 deb-package/DEBIAN/control
chmod 755 deb-package/var/lib/bitget-collector

# Build DEB package
dpkg-deb --build deb-package bitget-collector.deb

# Clean up
rm -rf deb-package collector
