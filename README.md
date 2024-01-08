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
efibootmgr -c -d /dev/nvme0n1 -p 1 -L "KLoader" -l "\EFI\KLoader\KLoader.efi" -u "url=data:text/html;base64,$(base64 --wrap=0 <(cat example.sh)) mount=root@UUID=4198c1da-48fc-44c3-888c-172facf208dc"
```

## Debug olddefconfig:
```
make olddefconfig
```
