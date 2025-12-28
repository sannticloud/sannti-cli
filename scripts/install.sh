#!/bin/bash
set -e

# Sannti CLI Installer
# Usage: curl -fsSL https://get.sannti.cloud/install.sh | bash

REPO="sannticloud/sannti-cli"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="sannti"
GITHUB_API="https://api.github.com/repos/${REPO}/releases/latest"
GITHUB_DOWNLOAD="https://github.com/${REPO}/releases/download"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${GREEN}ℹ${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

log_error() {
    echo -e "${RED}✗${NC} $1"
    exit 1
}

log_success() {
    echo -e "${GREEN}✓${NC} $1"
}

# Detect OS and Architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)

    case "$os" in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        *)
            log_error "Unsupported operating system: $os"
            ;;
    esac

    case "$arch" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            ;;
    esac

    log_info "Detected platform: ${OS}/${ARCH}"
}

# Get latest release version
get_latest_version() {
    log_info "Fetching latest version..."
    
    if command -v curl >/dev/null 2>&1; then
        VERSION=$(curl -s "${GITHUB_API}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        VERSION=$(wget -qO- "${GITHUB_API}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        log_error "curl or wget is required to download Sannti CLI"
    fi

    if [ -z "$VERSION" ]; then
        log_error "Failed to fetch latest version"
    fi

    log_info "Latest version: ${VERSION}"
}

# Download binary
download_binary() {
    local download_url="${GITHUB_DOWNLOAD}/${VERSION}/${BINARY_NAME}-${VERSION}-${OS}-${ARCH}"
    local tmp_file="/tmp/${BINARY_NAME}"

    log_info "Downloading Sannti CLI..."
    
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "${download_url}" -o "${tmp_file}"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "${download_url}" -O "${tmp_file}"
    fi

    if [ ! -f "${tmp_file}" ]; then
        log_error "Failed to download binary"
    fi

    chmod +x "${tmp_file}"
    log_success "Downloaded successfully"
}

# Install binary
install_binary() {
    log_info "Installing to ${INSTALL_DIR}..."

    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        log_warn "Requires sudo to install to ${INSTALL_DIR}"
        sudo mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    log_success "Installed to ${INSTALL_DIR}/${BINARY_NAME}"
}

# Verify installation
verify_installation() {
    if command -v ${BINARY_NAME} >/dev/null 2>&1; then
        local installed_version=$(${BINARY_NAME} version 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+')
        log_success "Sannti CLI ${installed_version} installed successfully!"
        return 0
    else
        log_error "Installation verification failed"
    fi
}

# Main installation flow
main() {
    echo "╔═══════════════════════════════════════════════════════╗"
    echo "║       Sannti CLI Installer                            ║"
    echo "║       Beyond Cloud. Without Barriers.                 ║"
    echo "╚═══════════════════════════════════════════════════════╝"
    echo ""

    detect_platform
    get_latest_version
    download_binary
    install_binary
    verify_installation

    echo ""
    echo "╔═══════════════════════════════════════════════════════╗"
    echo "║  Next steps:                                          ║"
    echo "║                                                       ║"
    echo "║  1. Configure your credentials:                      ║"
    echo "║     $ sannti configure                               ║"
    echo "║                                                       ║"
    echo "║  2. List available regions:                          ║"
    echo "║     $ sannti region list                             ║"
    echo "║                                                       ║"
    echo "║  3. View help:                                       ║"
    echo "║     $ sannti --help                                  ║"
    echo "║                                                       ║"
    echo "║  Documentation: https://docs.sannti.cloud            ║"
    echo "╚═══════════════════════════════════════════════════════╝"
}

# Run installer
main
