# Fullscreen Clock

A deliberately simple fullscreen clock written in Go using [Fyne](https://fyne.io/).

It is intended for a Linux kiosk/display where the only useful information is:

- the current time
- the current date in Dutch
- whether the system clock has been synchronised through NTP

When `timedatectl` reports that NTP is not synchronised, the clock displays:

    Tijd synchroniseren ...

The application has no window decorations and runs fullscreen. The mouse cursor is hidden.

## Requirements

For local development you need:

- Go
- GCC
- Linux graphics development libraries required by Fyne
- `timedatectl` from systemd

Fyne's Linux prerequisites include GCC and X11/Wayland development libraries. See the
[Fyne quick-start documentation](https://docs.fyne.io/started/quick/) for distribution-specific packages.

For Debian/Ubuntu:

```sh
sudo apt-get install \
    golang \
    gcc \
    libgl1-mesa-dev \
    xorg-dev \
    libwayland-dev \
    libxkbcommon-dev
```

For NixOS:
```sh
nix-shell -p go libGL pkg-config libX11.dev libxcursor libxi libxinerama libxrandr libxxf86vm libxkbcommon wayland
```


Update bundled assets:
```bash
go install fyne.io/tools/cmd/fyne@latest
fyne bundle -o bundled.go ChivoMono-ExtraBold.ttf
```
