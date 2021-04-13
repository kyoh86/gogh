#!/bin/bash
set -exuo pipefail

pacman -Syyu --noconfirm --needed
pacman -Syu --noconfirm --needed base-devel pacman-contrib

export VERSION="${GITHUB_REF##*/v}"
if [[ "${VERSION}" =~ ^([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
  :
else
  export VERSION="2.0.1"
fi
export COMMIT="${GITHUB_SHA}"
envsubst "\$VERSION \$COMMIT" < pkgbuild.template.sh > PKGBUILD
chmod 777 PKGBUILD
updpkgsums
cat PKGBUILD
