#!/bin/bash
# Sync KrankyBearClock to Ubuntu 18.04, compile, and retrieve binary
# Ubuntu 18.04 has glibc 2.27 which provides excellent compatibility

set -euo pipefail

# UPDATE THESE FOR YOUR UBUNTU 18.04 SYSTEM:
UBUNTU_USER="allan"
UBUNTU_HOST="192.168.1.18"  # Your Ubuntu 18.04 IP address
UBUNTU_PATH="/home/${UBUNTU_USER}/KrankyBearClock"

log() { printf '%s %s\n' "$(date '+%Y-%m-%d %H:%M:%S')" "$*"; }

run_remote() {
  local cmd="$1"
  ssh -o StrictHostKeyChecking=accept-new "${UBUNTU_USER}@${UBUNTU_HOST}" "set -euo pipefail; source ~/.bash_profile >/dev/null 2>&1 || true; cd ${UBUNTU_PATH}; ${cmd}"
}

echo "=== Syncing to Ubuntu 18.04 ==="
echo "Target: ${UBUNTU_USER}@${UBUNTU_HOST}:${UBUNTU_PATH}"

cp ReleaseNotes.txt Resources

rsync -av --delete \
  --exclude='bin/' \
  --exclude='.git/' \
  --exclude='.DS_Store' \
  --exclude='installers/' \
  --exclude='KrankyBearClock.app' \
  . "${UBUNTU_USER}@${UBUNTU_HOST}:${UBUNTU_PATH}/"

echo ""
echo "=== Compiling on Ubuntu 18.04 ==="
run_remote "./compile-linux.sh"

echo ""
echo "=== Retrieving compiled binary ==="
mkdir -p ./bin/
rm -f ./bin/KrankyBearClock-linux-*
if ! scp "${UBUNTU_USER}@${UBUNTU_HOST}:${UBUNTU_PATH}/bin/KrankyBearClock-linux-*" ./bin/
then
  echo "ERROR: Failed to copy compiled binaries back from Ubuntu." >&2
  exit 1
fi
ls -lh ./bin/KrankyBearClock-linux-* || true

echo ""
echo "=== Packaging on Ubuntu 18.04 ==="
run_remote "./package.sh linux"
run_remote "./package.sh linux ARCH=arm64"

echo "=== Retrieving installers ==="
mkdir -p ./installers/
rm -f ./installers/KrankyBearClock*.deb ./installers/KrankyBearClock*.rpm

if ! scp "${UBUNTU_USER}@${UBUNTU_HOST}:${UBUNTU_PATH}/installers/KrankyBearClock*.deb" ./installers/
then
  echo "WARNING: No .deb installer copied. Check remote build output." >&2
fi
if ! scp "${UBUNTU_USER}@${UBUNTU_HOST}:${UBUNTU_PATH}/installers/KrankyBearClock*.rpm" ./installers/
then
  echo "WARNING: No .rpm installer copied. Check remote build output." >&2
fi


# Prerequisites for Ubuntu 18.04:
# sudo apt update && sudo apt install -y gcc pkg-config libgl1-mesa-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev libx11-dev libxcursor-dev libxss-dev
# 
# Install Go 1.21+ on Ubuntu 18.04:
# wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
# sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
# echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
# source ~/.bashrc
#
# For KrankyBearClock, you'll also need (if not already installed):
# sudo apt install -y libgl1-mesa-dev libx11-dev libxcursor-dev libxinerama-dev libxi-dev libxrandr-dev libxss-dev libxxf86vm-dev


# "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
