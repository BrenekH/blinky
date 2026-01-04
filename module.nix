{ lib, pkgs, config, ... }:
with lib;
let
  # Shorter name to access final settings a user of this module HAS ACTUALLY SET.
  # cfg is a typical convention.
  cfg = config.services.blinky;

  daemonCommand = "${pkgs.blinky-daemon}/bin/blinkyd";
in {
  # Declare what settings a user of this module CAN SET.
  options.services.blinky = {
    enable = mkEnableOption "Blinky service";

    cli = {
      enable = mkEnableOption "Enable blinky command";
    };
  };

  # Define what other settings, services and resources should be active IF
  # a user of this module ENABLED this module
  # by setting "services.docker-compose-manager.enable = true;".
  config = mkIf cfg.enable {
    systemd.services.blinky-daemon = {
      serviceConfig = {
        ExecStart = daemonCommand;
        Restart = "on-failure";
        RestartSec = 5;
      };

      wantedBy = [ "default.target" ];
    };

    environment.systemPackages = [
      pkgs.blinky-daemon
    ] ++ (if cfg.cli.enable then [ pkgs.blinky-cli ] else [ ]);
  };
}
