
run-dev:
    #!/bin/sh
    set -e
    dir="tmp"
    mkdir -p "$dir"
    echo contents > "$dir/one"
    watchexec -e go --watch . --restart --debounce 1s  go run . -d "$dir"
    rm -rf "$dir"

build:
    go build -o fh .
