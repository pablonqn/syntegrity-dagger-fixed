#!/bin/bash

# Syntegrity Dagger Installer Script
# This script downloads the latest release of syntegrity-dagger

set -e

# Configuration
REPO="getsyntegrity/syntegrity-dagger"
BINARY_NAME="syntegrity-dagger"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
VERSION="${VERSION:-latest}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Detect OS and Architecture
detect_platform() {
    local os arch
    
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*)    os="windows" ;;
        MINGW*)     os="windows" ;;
        *)          log_error "Unsupported OS: $(uname -s)"; exit 1 ;;
    esac
    
    case "$(uname -m)" in
        x86_64)     arch="amd64" ;;
        arm64)      arch="arm64" ;;
        aarch64)    arch="arm64" ;;
        *)          log_error "Unsupported architecture: $(uname -m)"; exit 1 ;;
    esac
    
    echo "${os}-${arch}"
}

# Get latest release version
get_latest_version() {
    local version
    version=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [[ -z "$version" ]]; then
        log_error "Failed to get latest version"
        exit 1
    fi
    
    echo "$version"
}

# Download binary
download_binary() {
    local version=$1
    local platform=$2
    local os arch ext
    
    IFS='-' read -r os arch <<< "$platform"
    
    if [[ "$os" == "windows" ]]; then
        ext=".exe"
    else
        ext=""
    fi
    
    local filename="${BINARY_NAME}-${platform}${ext}"
    local url="https://github.com/${REPO}/releases/download/${version}/${filename}"
    local output_path="${INSTALL_DIR}/${BINARY_NAME}${ext}"
    
    log_info "Downloading ${BINARY_NAME} ${version} for ${platform}..."
    log_info "URL: ${url}"
    
    # Create install directory if it doesn't exist
    sudo mkdir -p "${INSTALL_DIR}"
    
    # Download the binary
    if curl -L -o "${output_path}" "${url}"; then
        log_success "Downloaded ${filename}"
    else
        log_error "Failed to download ${filename}"
        exit 1
    fi
    
    # Make it executable
    sudo chmod +x "${output_path}"
    
    # Verify installation
    if "${output_path}" --version >/dev/null 2>&1; then
        log_success "Installation successful!"
        log_info "Binary installed to: ${output_path}"
        log_info "Version: $("${output_path}" --version 2>/dev/null || echo "unknown")"
    else
        log_warning "Installation completed but version check failed"
    fi
}

# Main installation function
install() {
    local platform version
    
    log_info "Installing ${BINARY_NAME}..."
    
    # Detect platform
    platform=$(detect_platform)
    log_info "Detected platform: ${platform}"
    
    # Get version
    if [[ "$VERSION" == "latest" ]]; then
        version=$(get_latest_version)
        log_info "Latest version: ${version}"
    else
        version="$VERSION"
        log_info "Requested version: ${version}"
    fi
    
    # Download and install
    download_binary "$version" "$platform"
    
    log_success "Installation completed!"
    log_info "Usage: ${BINARY_NAME} --help"
}

# Show help
show_help() {
    cat << EOF
Syntegrity Dagger Installer

USAGE:
    $0 [OPTIONS]

OPTIONS:
    -v, --version VERSION    Install specific version (default: latest)
    -d, --dir DIRECTORY      Installation directory (default: /usr/local/bin)
    -h, --help              Show this help message

EXAMPLES:
    $0                      # Install latest version
    $0 -v v1.0.0           # Install specific version
    $0 -d ~/bin            # Install to custom directory

ENVIRONMENT VARIABLES:
    VERSION                 Version to install (default: latest)
    INSTALL_DIR            Installation directory (default: /usr/local/bin)

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -d|--dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Check dependencies
if ! command -v curl &> /dev/null; then
    log_error "curl is required but not installed"
    exit 1
fi

# Run installation
install
