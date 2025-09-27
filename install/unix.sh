#!/bin/bash

INSTALL_DIR="/usr/local/bin"
EXECUTABLE_NAME="ggh"
EXECUTABLE_PATH="$INSTALL_DIR/$EXECUTABLE_NAME"
USE_SUDO="false"
OS=""
ARCH=""

RED='\033[0;31m'
PURPLE='\033[0;35m'
GREEN='\033[0;32m'
NC='\033[0m'

setSystem() {
    ARCH="$(uname -m)"
    case "$ARCH" in
        i386|i686) ARCH="x86_64" ;;
        x86_64) ARCH="x86_64" ;;
        armv6*) ARCH="arm64" ;;
        armv7*) ARCH="arm64" ;;
        aarch64*) ARCH="arm64" ;;
    esac

    OS="$(uname | tr '[:upper:]' '[:lower:]')"
    if [ "$OS" = "linux" ]; then
        USE_SUDO="true"
    fi
    if [ "$OS" = "darwin" ] && uname -a | grep -qi "arm64"; then
        USE_SUDO="true"
    fi
}

runAsRoot() {
    local CMD=("$@")
    if [ "$USE_SUDO" = "true" ]; then
        printf "%bWe need sudo access to move GGH to %s%b\n" "$PURPLE" "$INSTALL_DIR" "$NC"
        CMD=(sudo "${CMD[@]}")
    fi
    "${CMD[@]}"
}

downloadBinary() {
    GITHUB_FILE="ggh_${OS}_${ARCH}"
    GITHUB_URL="https://github.com/byawitz/ggh/releases/latest/download/$GITHUB_FILE"
    curl -L --progress-bar --output "ggh-tmp" "$GITHUB_URL"
}

install() {
    if ! chmod +x "ggh-tmp"; then
        printf "%bFailed to set permissions ...%b\n" "$RED" "$NC"
        exit 1
    fi

    if ! runAsRoot mv "ggh-tmp" "$EXECUTABLE_PATH"; then
        printf "%bFailed to copy file ...%b\n" "$RED" "$NC"
        exit 1
    fi
}

printf "%bInstalling GGH%b\n\n" "$PURPLE" "$NC"

setSystem

downloadBinary

install

printf "\n%bGGH was installed successfully to:%b %s\n" "$GREEN" "$NC" "$EXECUTABLE_PATH"
