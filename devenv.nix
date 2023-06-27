{ pkgs, ... }:

{
  # https://devenv.sh/packages/
  packages = [ 
    pkgs.buf
    pkgs.go
    pkgs.git
    pkgs.git-cliff
    pkgs.govulncheck
    pkgs.python311
  ];
}
