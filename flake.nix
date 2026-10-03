{
  description = "go devshell and package, created by scaffolder";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          name = "go-devshell";

          packages = with pkgs; [
            go
            gopls
            gotools
            delve
            just
          ];
        };

        packages.lanyard = pkgs.buildGoModule {
          pname = "lanyard";
          version = "0.1.1";

          src = self;

          vendorHash = "sha256-UdRyylHSZ/b89cEArilaw6LcAw0epmpH6yCOKOiw9Gw=";

          subPackages = [ "." ];
          ldflags = [
            "-s"
            "-w"
          ];

          postInstall = ''
            mkdir -p $out/share/lanyard
            cp shell/lanyard.nu $out/share/lanyard/lanyard.nu
          '';

          meta = with pkgs.lib; {
            description = "Manager for ssh-agent and keys";
            license = licenses.unlicense;
            platforms = platforms.all;
          };
        };

        apps.lanyard = {
          type = "app";
          program = "${self.packages.${pkgs.stdenv.hostPlatform.system}.lanyard}/bin/lanyard";
        };
      }
    )
    // {
      homeModules.default = import ./home.nix { inherit self; };
    };
}
