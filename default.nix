let
        _pkgs = import <nixpkgs> {};
in
        { pkgs ? import (_pkgs.fetchFromGitHub {
                owner = "NixOS";
                repo = "nixpkgs-channels";
                rev = "9c048f4fb66adc33c6b379f2edefcb615fd53de6";
                sha256 = "18xbnfzj753bphzmgp74rn9is4n5ir4mvb4gp9lgpqrbfyy5dl2j";
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
