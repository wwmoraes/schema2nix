{
  pkgs,
  ...
}:

pkgs.buildGoApplication rec {
  pname = "schema2nix";
  version = "0.0.0";

  src =
    with pkgs.lib.fileset;
    toSource {
      root = ./.;
      fileset = unions [
        (fileFilter (file: file.hasExt "go") ./.)
        (maybeMissing ./go.sum)
        ./nix/module.gotmpl
        ./go.mod
      ];
    };

  modules = ./gomod2nix.toml;
  subPackages = [ "cmd/schema2nix" ];

  CGO_ENABLED = 0;
  GOFLAGS = "-trimpath";

  ldflags = [
    "-s"
    "-w"
    "-buildid="
    "-X main.version=${version}"
  ];

  meta = {
    description = "converts JSON schemas to nix options";
    homepage = "https://github.com/wwmoraes/schema2nix";
    license = pkgs.lib.licenses.mit;
    maintainers = [ pkgs.lib.maintainers.wwmoraes ];
    mainProgram = "schema2nix";
  };
}
