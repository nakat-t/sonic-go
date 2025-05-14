#!/bin/bash
set -euo pipefail

script_dir="$(dirname "$(realpath "$0")")"
source_dir="$script_dir/../submodules/sonic"
target_dir="$script_dir/../internal/cgosonic"

files=(
    "sonic.h"
    "sonic.c"
    "wave.h"
    "wave.c"
)

for file in "${files[@]}"; do
    src="$source_dir/$file"
    dest="$target_dir/$file"
    
    if [ -f "$src" ]; then
        cp "$src" "$dest"
    else
        echo "error: source file $src does not exist" >&2
        exit 1
    fi
done
