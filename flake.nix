{
  description = "Flake for installation and development of Blinky. Made with gomod2nix.";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  inputs.gomod2nix.inputs.flake-utils.follows = "flake-utils";

  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
    (flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          ### Modules ###
          nixosModules.default = self.nixosModules.${system}.blinky;
          nixosModules.blinky = ./module.nix;

          ## Overlays ##
          overlays.default = self.overlays.${system}.blinky;
          overlays.blinky = final: _prev: (pkgs: {
              blinkyd = self.packages.${system}.blinkyd;
              blinky = self.packages.${system}.blinky;
            }) final.pkgs;

          ### Packages ###
          packages.blinkyd = pkgs.callPackage ./. {
            inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
          };

          packages.blinky = pkgs.callPackage ./build-cli.nix {
            inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
          };

          packages.default = self.packages.${system}.blinkyd;

          ### Shells ###
          devShells.default = let
            goEnv = gomod2nix.legacyPackages.${system}.mkGoEnv { pwd = ./.; };
          in
            pkgs.mkShell {
              packages = [
                goEnv
                gomod2nix.legacyPackages.${system}.gomod2nix
              ];

              shellHook = ''
                echo "Updating gomod2nix.toml"
                ${gomod2nix.legacyPackages.${system}.gomod2nix}/bin/gomod2nix

                echo "Enabling dev environment"
                go version

                exec fish
              '';
            };
        })
    );
}
