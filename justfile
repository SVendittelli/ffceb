# See https://github.com/casey/just
set shell := ["zsh", "-uc"]
set windows-shell := ["pwsh.exe", "-NoLogo", "-Command"]

[private]
@default:
    just --list --justfile {{justfile()}}

export CGO_ENABLED := "1"

# go run with cgo enabled
run *args:
    go run main.go {{args}}

