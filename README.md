### Firestarter

## Linux wrapped in x86-64 EFI application for prepare and chainload to anything

Firestarter is minimized Linux especially designed to prepare the environment before chainload to anything.

## Build dependencies:

- docker buildx
- internet
- make

## Build:

```
make clean
make
```

## Bugs

- efivars FilePath separation is not supported (see https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#file-path-media-device-path)
- Crutched IPv4 DHCP bootfile-name option parsing logic with fake DHCP client
- IPv6 DHCP bootfile-url option parsing is not implemented
- IPv6 network stack and VLAN network layer are not tested at all
