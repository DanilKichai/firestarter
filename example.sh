#!/bin/bash -e

PREIFX="/mnt/root"
ROOT_UUID="$(findmnt --output UUID --noheadings "${PREIFX}")"

kexec \
  --load "${PREIFX}/boot/vmlinuz" \
  --initrd="${PREIFX}/boot/initrd.img" \
  --append="root=UUID=${ROOT_UUID} ro quiet splash vt.handoff=7" 
kexec --exec
