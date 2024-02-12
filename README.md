### KLoader

## EFI x86-64 Linux kernel loader with the power of bash scripting

KLoader is minimized Linux especially designed to prepare environment and chainload to another Linux kernel.

## Build dependencies:

- docker buildx
- git
- internet
- make

## Build:

```
make clean
make
```

## Install:

```
mkdir -p /boot/efi/EFI/KLoader
cp KLoader.efi /boot/efi/EFI/KLoader/
cp example/KLoader.sh /boot/
efibootmgr -c -d /dev/nvme0n1 -p 1 -L "KLoader" -l "\EFI\KLoader\KLoader.efi" -u "url=file:///mnt/root/boot/KLoader.sh mount=/mnt/root@/dev/disk/by-uuid/$(findmnt --output UUID --noheadings /)"
```
