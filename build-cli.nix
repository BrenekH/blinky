{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
, buildGoApplication ? pkgs.buildGoApplication
}:

buildGoApplication {
  pname = "blinky";
  name = "blinky";
  version = "0.1";
  pwd = ./.;
  src = ./.;
  subPackages = [ "cmd/blinky" ];
  modules = ./gomod2nix.toml;
  installPhase = ''
    install -Dm755 "$GOPATH/bin/blinky" "$out/bin/blinky"
  '';
}
