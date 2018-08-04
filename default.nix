
with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "ape";
  buildInputs = [
    pkgs.go_1_10
    pkgs.vndr
    pkgs.gnumake
    pkgs.gotools
    pkgs.golint
    pkgs.godef
    pkgs.gocode
  ];
}
