{
  description = "converts JSON schemas to nix options";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    flake-utils = {
      inputs.systems.follows = "systems";
      url = "github:numtide/flake-utils";
    };
    gomod2nix = {
      inputs.flake-utils.follows = "flake-utils";
      inputs.nixpkgs.follows = "nixpkgs";
      url = "github:tweag/gomod2nix";
    };
    nixpkgs.url = "github:NixOS/nixpkgs/25.11";
    nur = {
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-parts.follows = "flake-parts";
      url = "github:nix-community/NUR";
    };
    systems.url = "github:nix-systems/default";
    treefmt-nix = {
      inputs.nixpkgs.follows = "nixpkgs";
      url = "github:numtide/treefmt-nix";
    };
    unstable.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };
  nixConfig = {
    extra-substituters = [
      "https://wwmoraes.cachix.org/"
      "https://nix-community.cachix.org/"
    ];
    extra-trusted-public-keys = [
      "wwmoraes.cachix.org-1:N38Kgu19R66Jr62aX5rS466waVzT5p/Paq1g6uFFVyM="
      "nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs="
    ];
  };
  outputs =
    inputs@{
      self,
      ...
    }:
    (inputs.flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.treefmt-nix.flakeModule
      ];

      flake = {
        overlays = {
          default =
            final: prev:
            let
              inherit (self.packages.${prev.stdenv.hostPlatform.system}) schema2nix;
            in
            {
              inherit schema2nix;
              mkOptionFromSchema =
                {
                  src,
                  name ? "${baseNameOf src}.nix",
                }:
                prev.stdenv.mkDerivation {
                  inherit src name;
                  buildInputs = [
                    schema2nix
                  ];
                  buildCommand = ''
                    schema2nix "$src" > "$out"
                  '';
                };
            };
        };
      };

      perSystem =
        {
          pkgs,
          self',
          system,
          ...
        }:
        {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;
            overlays = [
              inputs.gomod2nix.overlays.default
              inputs.nur.overlays.default
              self.overlays.default
              ## add wwmoraes maintainer
              (final: prev: {
                lib = prev.lib.extend (
                  final: prev: {
                    maintainers = prev.maintainers // {
                      wwmoraes = {
                        email = "nixpkgs@artero.dev";
                        github = "wwmoraes";
                        githubId = 682095;
                        keys = [ { fingerprint = "32B4 330B 1B66 828E 4A96  9EEB EED9 9464 5D7C 9BDE"; } ];
                        matrix = "@wwmoraes:hachyderm.io";
                        name = "William Artero";
                      };
                    };
                  }
                );
              })
              (final: prev: {
                unstable = import inputs.unstable { inherit system; };
              })
            ];
            config = { };
          };
          devShells = import ./shells.nix (pkgs // { inherit pkgs; });
          packages = {
            default = self'.packages.schema2nix;
            schema2nix = import ./default.nix { inherit pkgs; };
          };
        };
      systems = import inputs.systems;
    });
}
