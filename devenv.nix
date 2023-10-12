{ pkgs, ... }:

{
  # https://devenv.sh/packages/
  packages = [ 
    pkgs.buf
    pkgs.go_1_21
    pkgs.git
    pkgs.git-cliff
    pkgs.govulncheck
    pkgs.gopls
    pkgs.golint
    pkgs.mysql-shell
    pkgs.postgresql_15
    pkgs.python311
  ];
}
