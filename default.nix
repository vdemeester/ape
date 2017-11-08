let
        _pkgs = import <nixpkgs> {};
in
        { pkgs ? import (_pkgs.fetchFromGitHub {
                owner = "NixOS";
                repo = "nixpkgs-channels";
                rev = "d0d905668c010b65795b57afdf7f0360aac6245b";
                sha256 = "1kqxfmsik1s1jsmim20n5l4kq6wq8743h5h17igfxxbbwwqry88l";
        }) {}
}:

pkgs.stdenv.mkDerivation rec {
        name = "go-projects";
        env = pkgs.buildEnv { name = name; paths = buildInputs; };
        buildInputs = [
                pkgs.go_1_9
                pkgs.vndr
                pkgs.gnumake
                pkgs.gotools
                pkgs.golint
                pkgs.godef
                pkgs.gocode
        ];
}
