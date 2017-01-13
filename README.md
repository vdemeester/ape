# 🐒 ape [![Build Status](https://travis-ci.org/vdemeester/ape.svg?branch=master)](https://travis-ci.org/vdemeester/ape)

`ape` is git mirror updater, nothing more. It reads a simple file that holds
a list of git repository (URL) and the upstream to rebase against.

```bash
λ cat $HOME/.config/ape.conf
git@github.com:vdemeester/docker.git https://github.com/docker/docker.git
git@github.com:vdemeester/libcompose.git https://github.com/docker/libcompose.git
# […]
git@github.com:vdemeester/nixpkgs.git https://github.com/NixOS/nixpkgs.git
# […]
git@github.com:vdemeester/traefik.git https://github.com/containous/traefik.git
# […]
λ ape up ~/var/mirrors
🐒 docker
🐒 libcompose..
🙈 cloning git@github.com:vdemeester/docker.git
🙈 cloning git@github.com:vdemeester/libcompose.git
🐒 traefik..
🙈 cloning git@github.com:vdemeester/traefik.git
🐒 nixpkgs..
🙈 cloning git@github.com:vdemeester/nixpkgs.git
🙉 add upstream https://github.com/containous/traefik.git
🙉 add upstream https://github.com/docker/libcompose.git
🙊 fetch and rebase libcompose
🙊 fetch and rebase traefik
🙉 add upstream https://github.com/docker/docker.git
🙉 add upstream https://github.com/NixOS/nixpkgs.git
🐵 push to origin libcompose
🐵 push to origin traefik
🙊 fetch and rebase nixpkgs
🙊 fetch and rebase docker
🐵 push to origin docker
🐵 push to origin nixpkgs
# […] Later that day
λ ape up ~/var/mirrors
🐒 docker
🐒 libcompose..
🐒 traefik..
🐒 nixpkgs..
🙊 fetch and rebase libcompose
🙊 fetch and rebase traefik
🐵 push to origin libcompose
🐵 push to origin traefik
🙊 fetch and rebase nixpkgs
🙊 fetch and rebase docker
🐵 push to origin docker
🐵 push to origin nixpkgs
```