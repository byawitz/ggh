#!/bin/bash

INSTALL_DIR="/usr/local/bin"
EXECUTABLE_NAME=ggh
EXECUTABLE_PATH="$INSTALL_DIR/$EXECUTABLE_NAME"
USE_SUDO="false"
OS=""
ARCH=""

RED='\033[0;31m'
PURPLE='\033[0;35m'
GREEN='\033[0;32m'
NC='\033[0m'

setSystem() {
    ARCH=$(uname -m)
    case $ARCH in
        i386|i686) ARCH="x86_64" ;;
        x86_64) ARCH="x86_64";;
        armv6*) ARCH="arm64" ;;
        armv7*) ARCH="arm64" ;;
        aarch64*) ARCH="arm64" ;;
    esac

    OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')
    if [ "$OS" = "linux" ]; then
        USE_SUDO="true"
    fi
    if [ "$OS" = "darwin" ] && [[ "$(uname -a)" = *ARM64* ]]; then
        USE_SUDO="true"
    fi
}

runAsRoot() {
    local CMD="$*"
        if [ "$USE_SUDO" = "true" ]; then
          printf "${PURPLE}We need sudo access to add mv GGH to $INSTALL_DIR ${NC}\n"
          CMD="sudo $CMD"
        fi
    $CMD
}

downloadBinary() {
    GITHUB_FILE="ggh_${OS}_${ARCH}"
    GITHUB_URL="https://github.com/byawitz/ggh/releases/latest/download/$GITHUB_FILE"
    curl $GITHUB_URL --location --progress-bar --output "ggh-tmp"

}

install() {
    chmod +x "ggh-tmp"
    if [ $? -ne 0 ]; then
        printf "${RED}Failed to set permissions ... ${NC}"
        exit 1
    fi

    runAsRoot mv "ggh-tmp" "$EXECUTABLE_PATH"
    if [ $? -ne 0 ]; then
        printf "${RED}Failed to copy file ... ${NC}"
        exit 1
    fi
}



printf "${PURPLE}Installing GGH ${NC}\n\n"

setSystem
downloadBinary
install
printf "\n${GREEN}GGH was installed successfully to:${NC} $EXECUTABLE_PATH\n"