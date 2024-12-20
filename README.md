### ArchShell

ArchShell is minimized Linux wrapped in the EFI x86-64 application.

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

- EFI FilePath separation is not supported (see https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#file-path-media-device-path)
- IPv6 network stack and VLAN network layer are not tested at all
