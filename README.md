# BusyBox Loader
BusyBox Loader is the customizable linux kernel loader for x86_64-efi systems with helps of BusyBox shell scripting.

#### Build dependencies:
- internet connection
- docker with buildx plugin
- make tool

#### Build:
- `make`
- `make all`
- `make create_builder`
- `make build_bbloader`
- `make remove_bbloader`
- `make remove_builder`
- `make clean`

#### Usage:
- `efibootmgr -c -d /dev/nvme0n1 -p 1 -L "BBLoader" -l '\bbloader.efi'`
- `efibootmgr -c -d /dev/nvme0n1 -p 1 -L "BBLoader" -l '\bbloader.efi' -u 'source=H4sIAAAAAAACA0tNzshX8EjNyclXKM8vyklR5AIA3SQIaBIAAAA='`, where source is base64 encoded gzipped BusyBox shell script source
- `efibootmgr -c -d /dev/nvme0n1 -p 1 -L "BBLoader" -l '\bbloader.efi' -u 'initrd=\initramfs.override'`
