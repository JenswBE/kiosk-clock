# AGENTS.md

This file contains useful hints for working with this project.

## Project Structure

- `cmd/clock` — Application entrypoint
- `internal/config` — Configuration file handling
- `internal/system` — System functionality (backlight, WiFi/NTP, power management)
- `internal/ui` — User interface components
- `internal/assets` — Bundled fonts and assets

## Important Notes

- This project uses Go 1.26
- The main framework is Fyne v2
- The project uses Nix for dependency management
- Configuration of Fyne is handled via TOML files
- Power management features require Linux system access
- When updating any Golang dependency, rebuild with `nix build .#default --print-build-logs` and update `vendorHash` in `flake.nix` with the hash from the build output
