let
  nixpkgs = import ./pkgs.nix;
  pkgs = import nixpkgs {};
  drv = pkgs.callPackage ./derivation.nix {};
in 
  drv.overrideAttrs (attrs: {
    src = null;
    nativeBuildInputs = [ pkgs.operator-sdk pkgs.git ] ++ attrs.nativeBuildInputs;
    shellHook = ''
      eval "$configurePhase"
    ''; 
  })
