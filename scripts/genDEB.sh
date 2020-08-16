#!/bin/bash

VERSION_STRING=1.0.0
DEB_PACKAGE_DESCRIPTION="Darn simple update server"

mkdir ../.dpkg
mkdir -p ../.dpkg/systemd/system
cp ./dsus.service ../.dpkg/systemd/system/
mkdir -p ../.dpkg/dsus/certs

if which fpm; then
    fpm -s dir \
      -t deb \
      --name dsus \
      --version $VERSION_STRING \
      --description '${DEB_PACKAGE_DESCRIPTION}' \
      -p dsus-${VERSION_STRING}.deb \
      --depends curl \
      --depends openssl \
      --after-install ./afterInstall.sh \
      ../dsus=/usr/bin/ \
      ../.dpkg/=/etc/
    rm -rf ../.dpkg
else
    echo "fpm not installed or not reachable"
    exit 1
fi