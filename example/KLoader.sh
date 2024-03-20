#!/usr/bin/env bash

ROOT="/dev/disk/by-uuid/98ee92a2-d538-44aa-bfd5-4d8d8f8896f0"

udevadm \
  wait --timeout=30 \
    "${ROOT}"
mount \
  --options ro \
  --source "${ROOT}" \
  --target /mnt

kexec \
  --load "/mnt/boot/vmlinuz" \
  --initrd="/mnt/boot/initrd.img" \
  --append="root=${ROOT} ro" 
kexec \
  --exec
