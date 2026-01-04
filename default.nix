{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  ),
  buildGoApplication ? pkgs.buildGoApplication,
}:

buildGoApplication {
  pname = "blinkyd";
  name = "blinkyd";
  version = "0.1";
  pwd = ./.;
  src = ./.;
  subPackages = [ "cmd/blinkyd" ];
  modules = ./gomod2nix.toml;

  # buildInputs = runtime dependencies
  # nativeBuildInputs = build time dependencies

  installPhase = ''
    install -Dm755 "$GOPATH/bin/blinkyd" "$out/bin/blinkyd"
  '';
}
