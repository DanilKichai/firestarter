### KLoader

## EFI x86-64 Linux kernel loader with the power of bash scripting

KLoader is minimized Linux especially designed to prepare environment and chainload to another Linux kernel.

## Build dependencies:

- docker buildx
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
efibootmgr -c -d /dev/nvme0n1 -p 1 -L "KLoader" -l "\EFI\KLoader\KLoader.efi"
```

## Configure:
```
sed --expression="s/98ee92a2-d538-44aa-bfd5-4d8d8f8896f0/$(findmnt --output UUID --noheadings /)/" <example/KLoader.sh
```
