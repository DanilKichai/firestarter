# Kexec Loader
Kexec Loader is the customizable linux kernel loader for x86_64-efi systems with helps of shell scripting.

#### Build dependencies:
- internet connection
- docker with buildx plugin
- make tool

#### Build:
- make clean
- make olddefconfigconfig
- edit linux.conf
- make

#### Usage:
- efibootmgr -c -d /dev/nvme0n1 -p 1 -L "kloader" -l '\kloader.efi'
- efibootmgr -c -d /dev/nvme0n1 -p 1 -L "kloader" -l '\kloader.efi' -u 'source=H4sIAAAAAAACA0tNzshX8EjNyclXKM8vyklR5AIA3SQIaBIAAAA='
- efibootmgr -c -d /dev/nvme0n1 -p 1 -L "kloader" -l '\kloader.efi' -u 'initrd=\initramfs.override'
