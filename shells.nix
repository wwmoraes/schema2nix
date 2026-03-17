{
  pkgs,
  ...
}:
rec {
  default = pkgs.mkShell {
    nativeBuildInputs = [
      # keep-sorted start
      (pkgs.mkGoEnv { pwd = ./.; })
      pkgs.cocogitto
      pkgs.goreleaser
      pkgs.remake
      pkgs.unstable.go
      pkgs.unstable.golangci-lint
      # keep-sorted end
    ];
  };

  ci = default.overrideAttrs (
    final: prev: {
      nativeBuildInputs = [
        # keep-sorted start
        # keep-sorted end
      ]
      ++ prev.nativeBuildInputs;

      shellHook = ''
        export GOCACHE=$(go env GOCACHE)
        export GOMODCACHE=$(go env GOMODCACHE)
      '';
    }
  );

  terminal = default.overrideAttrs (
    final: prev: {
      nativeBuildInputs = [
        # keep-sorted start
        pkgs.gomod2nix
        pkgs.nix-update
        pkgs.unstable.gotests
        pkgs.unstable.gotools
        # keep-sorted end
      ]
      ++ prev.nativeBuildInputs;

      shellHook = ''
        cog install-hook --all --overwrite
      '';
    }
  );
}
