{
  description = "A simple fullscreen clock";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachSystem
      [ "x86_64-linux" ]
      (system:
        let
          pkgs = import nixpkgs {
            inherit system;
          };
        in
        {
          packages.default = pkgs.buildGoModule {
            pname = "clock";
            version = "0.1.0";

            src = ./.;

            # Set this to the hash reported by Nix if dependencies need
            # to be fetched. null means don't use vendored dependencies.
            vendorHash = "sha256-veJbO0MXmDwfI7IyKoDTy2/jaxUy4whgzE/BExqmIbA=";

            subPackages = [ "cmd/clock" ];

            nativeBuildInputs = with pkgs; [
              pkg-config
              makeWrapper
            ];

            buildInputs = with pkgs; [
              # OpenGL
              libGL

              # X11
              libX11
              libxcursor
              libxi
              libxinerama
              libxrandr
              libxxf86vm

              # Keyboard/input
              libxkbcommon

              # Wayland
              wayland
            ];

            postInstall = ''
              wrapProgram "$out/bin/clock" \
                --prefix PATH : ${pkgs.lib.makeBinPath [
                  pkgs.brightnessctl
                ]}
            '';

            ldflags = [
              "-s"
              "-w"
            ];

            meta = {
              mainProgram = "clock";
              description = "A simple fullscreen clock";
              platforms = [ "x86_64-linux" ];
            };
          };
        }
      );
}
