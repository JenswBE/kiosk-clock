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
CLOCK_VERSION=v0.1.1
CLOCK_VERSION_MSG="Fixed vendor hash"

# Ensure cachix is installed
nix-shell -p cachix

# Go to https://app.cachix.org/cache/jenswbe/settings/authtokens and create a new read/write token
cachix authtoken AUTH_TOKEN

# Build flake
path=$(nix build --no-link --print-out-paths)

# Push to Cachix
echo "$path" | cachix push jenswbe

# Pin version
cachix pin jenswbe "${CLOCK_VERSION:?}" "$path"

# Git tag
git commit ...
git tag -a -m "${CLOCK_VERSION_MSG:?}" "${CLOCK_VERSION:?}"
```

## Miscellaneous

### Add bundled assets

```bash
go install fyne.io/tools/cmd/fyne@latest
fyne bundle -o bundled.go ChivoMono-ExtraBold.ttf
```
