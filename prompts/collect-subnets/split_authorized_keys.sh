#!/usr/bin/env bash

set -euo pipefail

INPUT="${1:-authorized_keys}"
OUTDIR="${2:-authorized_keys.d}"

mkdir -p "$OUTDIR"

i=0
while IFS= read -r line || [[ -n "$line" ]]; do
  [[ -z "$line" ]] && continue
  i=$((i+1))
  printf '%s\n' "$line" > "$OUTDIR/key_$(printf '%03d' "$i").pub"
done < "$INPUT"

echo "Written $i files to $OUTDIR"
