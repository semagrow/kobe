{ buildGoModule }:

buildGoModule {
  pname = "kobe-operator";
  version = "0.0.1";
  src = ./.;
  vendorSha256 = null;
}
