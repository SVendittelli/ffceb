# See https://github.com/casey/just
set shell := ["zsh", "-uc"]
set windows-shell := ["pwsh.exe", "-NoLogo", "-Command"]

[private]
@default:
    just --list --justfile {{justfile()}}

run args $CGO_ENABLED="1":
    go run main.go {{args}}

