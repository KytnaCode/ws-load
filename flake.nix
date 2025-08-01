{
  description = "A very basic Go flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    systems.url = "github:nix-systems/default";
  };

  outputs = {
    self,
    nixpkgs,
    systems,
  }: let
    eachSystem = nixpkgs.lib.genAttrs (import systems);
  in {
    devShells = eachSystem (system: let
      pkgs = import nixpkgs {inherit system;};
    in {
      default = pkgs.mkShell {
        buildInputs = with pkgs; [
          git

          # Go packages
          cobra-cli
          go
          golangci-lint
        ];
      };
    });

    packages = eachSystem (system: let
      pkgs = import nixpkgs {inherit system;};
      inherit (pkgs) lib;
    in rec {
      ws-load = pkgs.buildGoModule (finalAttrs: {
        pname = "ws-load";
        version = "0.1.0";

        src = ./.;

        vendorHash = "sha256-7Xjo1tAlHtzPnnbjQMHdo7Hae8HbJwdQCa+YLT/3T/s=";

        meta = {
          description = "websocket load generator written in go.";
          longDescription = ''
            ws-load is a simple websocket load generator with an easy-to-use interface.
          '';
          mainProgram = "ws-load";
          platforms = import systems;
          homepage = "https://github.com/kytnacode/ws-load";
          license = lib.licenses.mit;
        };
      });

      default = ws-load;
    });
  };
}
