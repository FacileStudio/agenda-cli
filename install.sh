#!/usr/bin/env bash
set -euo pipefail

if ! command -v tar >/dev/null 2>&1; then
  echo "tar is required" >&2
  exit 1
fi

if [ -z "${AGENDA_CLI_VERSION:-}" ]; then
  AGENDA_CLI_VERSION="$(curl -fsSL https://api.github.com/repos/FacileStudio/agenda-cli/releases/latest | grep tag_name | cut -d '"' -f4)"
fi

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "unsupported arch: $ARCH" >&2; exit 1 ;;
esac

case "$OS" in
  linux|darwin) ;;
  *) echo "unsupported os: $OS" >&2; exit 1 ;;
esac

URL="https://github.com/FacileStudio/agenda-cli/releases/download/${AGENDA_CLI_VERSION}/agenda_${AGENDA_CLI_VERSION#v}_${OS}_${ARCH}.tar.gz"
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

curl -fsSL "$URL" -o "$TMPDIR/agenda.tar.gz"
tar -xzf "$TMPDIR/agenda.tar.gz" -C "$TMPDIR"

sudo mv "$TMPDIR/agenda" /usr/local/bin/agenda
sudo chmod +x /usr/local/bin/agenda

echo "agenda-cli installed"
