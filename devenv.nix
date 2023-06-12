{ pkgs, ... }:

{
  # https://devenv.sh/packages/
  packages = [ 
    pkgs.buf
    pkgs.go
    pkgs.git
    pkgs.git-cliff
  ];
}
