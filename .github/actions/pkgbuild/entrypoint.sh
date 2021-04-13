#!/bin/bash
set -euo pipefail

pacman -Syu --noconfirm --needed base-devel pacman-contrib

# Makepkg does not allow running as root
# Create a new user `builder`
# `builder` needs to have a home directory because some PKGBUILDs will try to
# write to it (e.g. for cache)
useradd builder -m

# When installing dependencies, makepkg will use sudo
# Give user `builder` passwordless sudo access
echo "builder ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# Give all users (particularly builder) full access to these files
chmod -R a+rw .

export VERSION="${GITHUB_REF##*/v}"
if [[ "${VERSION}" =~ ^([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
  :
else
  export VERSION="2.0.1"
fi
export COMMIT="${GITHUB_SHA}"
export NAME="${GITHUB_REPOSITORY#*/}"
export OWNER="${GITHUB_REPOSITORY%/*}"
export REPOSITORY="${GITHUB_REPOSITORY}"
export SERVER_URL="${GITHUB_SERVER_URL}"
if [ -n "${EMAIL}" ]; then
  export EMAIL=" <${EMAIL}>"
fi
export DESCRIPTION
export LICENSE
envsubst "\$VERSION \$COMMIT \$NAME \$DESCRIPTION \$LICENSE \$OWNER \$REPOSITORY \$SERVER_URL \$EMAIL" < PKGBUILD.tmpl > PKGBUILD
chmod 777 PKGBUILD
sudo -H -u builder updpkgsums PKGBUILD
sudo -H -u builder makepkg --printsrcinfo PKGBUILD | sudo -H -u builder tee .SRCINFO
