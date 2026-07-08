#!/usr/bin/env bash
set -euo pipefail

REPO="bniladridas/kit"
INSTALL_DIR="${KIT_INSTALL_DIR:-$HOME/.local/bin}"
BIN_NAME="kit"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
  printf "${GREEN}[INFO]${NC} %s\n" "$1"
}

warn() {
  printf "${YELLOW}[WARN]${NC} %s\n" "$1"
}

error() {
  printf "${RED}[ERROR]${NC} %s\n" "$1" >&2
}

# Detect OS
detect_os() {
  case "$(uname -s)" in
    Linux*)  echo "linux" ;;
    Darwin*) echo "darwin" ;;
    *)       error "Unsupported OS: $(uname -s)"; exit 1 ;;
  esac
}

# Detect architecture
detect_arch() {
  case "$(uname -m)" in
    x86_64)  echo "amd64" ;;
    arm64)   echo "arm64" ;;
    aarch64) echo "arm64" ;;
    *)       error "Unsupported architecture: $(uname -m)"; exit 1 ;;
  esac
}

# Get latest release tag from GitHub API
get_latest_version() {
  local version
  version=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
  if [[ -z "$version" ]]; then
    error "Could not determine latest version"
    exit 1
  fi
  echo "$version"
}

# Download and verify checksum
download_with_checksum() {
  local url="$1"
  local checksum_url="$2"
  local output="$3"

  info "Downloading ${BIN_NAME}..."
  curl -fsSL "$url" -o "$output"

  if [[ ! -f "$output" ]]; then
    error "Download failed"
    exit 1
  fi

  info "Verifying checksum..."
  local expected_checksum
  expected_checksum=$(curl -fsSL "$checksum_url" | awk '{print $1}')

  if [[ -z "$expected_checksum" ]]; then
    warn "Could not fetch checksum, skipping verification"
    return 0
  fi

  local actual_checksum
  actual_checksum=$(sha256sum "$output" | awk '{print $1}')

  if [[ "$actual_checksum" != "$expected_checksum" ]]; then
    error "Checksum verification failed"
    error "Expected: $expected_checksum"
    error "Actual:   $actual_checksum"
    rm -f "$output"
    exit 1
  fi

  info "Checksum verified"
}

# Install binary
install_binary() {
  local tmpdir="$1"
  local binary="$2"

  info "Installing ${BIN_NAME} to ${INSTALL_DIR}"

  if [[ ! -d "$INSTALL_DIR" ]]; then
    mkdir -p "$INSTALL_DIR"
  fi

  cp "$binary" "${INSTALL_DIR}/${BIN_NAME}"
  chmod +x "${INSTALL_DIR}/${BIN_NAME}"

  # Ensure install dir is in PATH
  if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
    warn "${INSTALL_DIR} is not in your PATH"
    warn "Add this to your shell config:"
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
  fi

  info "Successfully installed ${BIN_NAME} to ${INSTALL_DIR}/${BIN_NAME}"
}

# Main
main() {
  info "Installing Kit CLI"

  local os arch version
  os=$(detect_os)
  arch=$(detect_arch)
  version=$(get_latest_version)

  info "Detected OS: ${os}, Architecture: ${arch}"
  info "Latest version: ${version}"

  local archive_name="${BIN_NAME}-${version}-${os}-${arch}"
  local archive_ext="tar.gz"
  if [[ "$os" == "windows" ]]; then
    archive_ext="zip"
  fi

  local archive_file="${archive_name}.${archive_ext}"
  local download_url="https://github.com/${REPO}/releases/download/${version}/${archive_file}"
  local checksum_url="https://github.com/${REPO}/releases/download/${version}/checksums.txt"

  local tmpdir
  tmpdir=$(mktemp -d)
  trap 'rm -rf "$tmpdir"' EXIT

  local archive_path="${tmpdir}/${archive_file}"
  local binary_path="${tmpdir}/${BIN_NAME}"

  download_with_checksum "$download_url" "$checksum_url" "$archive_path"

  info "Extracting..."
  if [[ "$archive_ext" == "tar.gz" ]]; then
    tar -xzf "$archive_path" -C "$tmpdir" "$archive_name/$BIN_NAME" || tar -xzf "$archive_path" -C "$tmpdir" "$BIN_NAME"
  else
    unzip -q "$archive_path" -d "$tmpdir"
  fi

  # Find the binary in the extracted files
  binary_path=$(find "$tmpdir" -name "$BIN_NAME" -type f | head -1)

  if [[ -z "$binary_path" || ! -f "$binary_path" ]]; then
    error "Binary not found in archive"
    exit 1
  fi

  install_binary "$tmpdir" "$binary_path"

  info "Installation complete!"
  info "Run '${BIN_NAME} --help' to get started"
}

main "$@"
