#!/bin/bash -e

PREIFX="/mnt/root"
UUID="$(findmnt --output UUID --noheadings "${PREIFX}")"

kexec \
  --load "${PREIFX}/boot/vmlinuz" \
  --initrd="${PREIFX}/boot/initrd.img" \
  --append="root=UUID=${UUID} ro quiet splash vt.handoff=7" 
kexec --exec
