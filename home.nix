{
  self,
}:
{
  config,
  lib,
  pkgs,
  ...
}:

let
  cfg = config.programs.lanyard;
in
{

  options.programs.lanyard = {
    enable = lib.mkEnableOption "Enable lanyard";
    package = lib.mkOption {
      type = lib.types.package;
      default = self.packages.${pkgs.stdenv.hostPlatform.system}.lanyard;
      description = "The lanyard package to use";
    };
    enableNushellIntegration = lib.mkEnableOption "Enable Nushell integration";
    keys = lib.mkOption {
      type = lib.types.lines;
      default = "";
      example = ''
        id_ed25519
      '';
      description = "List of SSH keys for lanyard to manage; stored in ~/.ssh/";
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [ cfg.package ];

    programs.nushell = lib.mkIf cfg.enableNushellIntegration {
      extraConfig = ''
        source ${cfg.package}/share/lanyard/lanyard.nu
      '';
    };

    xdg.configFile."lanyard/keys.conf" = lib.mkIf (cfg.keys != "") {
      text = cfg.keys;
    };
  };

}
