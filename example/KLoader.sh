#!/usr/bin/env bash

PREIFX="/mnt/root"
UUID="$(findmnt --output UUID --noheadings "${PREIFX}")"
APPEND="ro"

kexec \
  --load "${PREIFX}/boot/vmlinuz" \
  --initrd="${PREIFX}/boot/initrd.img" \
  --append="root=/dev/disk/by-uuid/${UUID} ${APPEND}" 
kexec --exec
