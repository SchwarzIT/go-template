{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:nix-community/gomod2nix";
  };

  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        name = "go-template";
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        };
      in
      rec {
        packages = {
          default = packages.${name};
          ${name} = pkgs.buildGoApplication {
            pname = name;
            version = pkgs.lib.removeSuffix "\n" (builtins.readFile ./config/version.txt);
            src = ./.;
            modules = ./nix/gomod2nix.toml;
            subPackages = [ "cmd/gt" ];
            ldflags = [ "-s" "-w" ];
            meta = with pkgs.lib; {
              description = "go/template is a tool for jumpstarting production-ready Golang projects quickly.";
              homepage = "https://github.com/schwarzit/go-template";
              license = licenses.asl20; # Apache License 2.0
            };
          };
        };

        apps = {
          default = flake-utils.lib.mkApp {
            drv = packages.default;
            exePath = "/bin/gt";
          };
        };
      });
}
