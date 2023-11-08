### KLoader

## EFI x86-64 Linux kernel loader with the power of bash scripting

KLoader is minimized Linux especially designed to prepare environment and chainload to another Linux kernel.

## Build dependencies:

- internet connection
- docker buildx
- GNU make

## Build:

```
make clean
make
```

## Prepare source argument:

```
gzip --best < source.sh | base64 --wrap=0
```

## Install:

```
efibootmgr -c -d /dev/nvme0n1 -p 1 -L "KLoader" -l '\EFI\KLoader\KLoader.efi' -u 'source=H4sIAAAAAAACA0tNzshX8EjNyclXKM8vyklR5AIA3SQIaBIAAAA='
```

## Debug olddefconfig:
```
make olddefconfig
```
