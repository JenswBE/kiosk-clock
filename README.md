# Fullscreen Clock

A simple fullscreen clock written in Go using [Fyne](https://fyne.io/).

## Local development

```bash
# Ensure deps are installed
nix-shell -p go libGL pkg-config libX11.dev libxcursor libxi libxinerama libxrandr libxxf86vm libxkbcommon wayland brightnessctl

# Run go tests and fixes
go mod tidy
go test ./...
go fix ./...

# Validate flake
nix flake check

# Flake test build
nix build .#default --print-build-logs
```

## Publishing new version

```bash
# Settings
CLOCK_VERSION=v0.1.2
CLOCK_VERSION_MSG="Bump dependencies"

# Ensure Golang is mod is clean
go mod tidy
go test ./...
go fix ./...

# Discard vendorHash
sed -i -E 's/(\s*)vendorHash.*/\1vendorHash = ""/' flake.nix

# Build to get new hash => Copy/paste hash into flake.nix
nix build .#default --print-build-logs

# Validate build with updated vendorHash
nix build .#default --print-build-logs

# Git tag
git add -A
git commit -m "${CLOCK_VERSION_MSG:?}"
git push origin main
git tag -a -m "${CLOCK_VERSION_MSG:?}" "${CLOCK_VERSION:?}"
git push --tags origin main

# Ensure cachix is installed
nix-shell -p cachix

# Go to https://app.cachix.org/cache/jenswbe/settings/authtokens and create a new read/write token
cachix authtoken AUTH_TOKEN

# Build flake
CLOCK_PATH=$(nix build --no-link --print-out-paths "github:jenswbe/kiosk-clock?ref=${CLOCK_VERSION:?}")

# Push to Cachix
echo "${CLOCK_PATH:?}" | cachix push jenswbe

# Pin version
cachix pin jenswbe "${CLOCK_VERSION:?}" "${CLOCK_PATH:?}"
```

## Miscellaneous

### Add bundled assets

```bash
go install fyne.io/tools/cmd/fyne@latest
fyne bundle -o bundled.go ChivoMono-ExtraBold.ttf
```
