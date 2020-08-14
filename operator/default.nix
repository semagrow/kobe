let 
  nixpkgs = import ./pkgs.nix;
  pkgs = import nixpkgs {
    config = { };
    overlays = [
      (import ./overlay.nix)
    ];
  };
in pkgs.kobe.operator
