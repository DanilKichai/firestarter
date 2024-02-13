# syntax=docker/dockerfile:1.2.1

FROM --platform=linux/amd64 archlinux:base as builder
  RUN \
    pacman \
      --sync \
      --sysupgrade \
      --refresh \
      --noconfirm \
      --needed \
      core/linux
  ADD toolchain/ /

  RUN \
    mkdir \
      --parents \
      /tmp/pacman \
      /target && \
    pacman \
      --root /target \
      --dbpath /tmp/pacman \
      --sync \
      --sysupgrade \
      --refresh \
      --noconfirm \
      --needed \
      core/bash \
      core/curl \
      core/coreutils \
      core/cryptsetup \
      core/dosfstools \
      core/e2fsprogs \
      core/efibootmgr \
      core/gawk \
      core/grep \
      core/gzip \
      core/iproute2 \
      core/iputils \
      core/less \
      core/links \
      core/lvm2 \
      core/nano \
      core/sed \
      core/traceroute \
      core/systemd \
      core/udev \
      core/util-linux \
      extra/clevis \
      extra/kexec-tools \
      extra/qrencode \
      extra/sbsigntools \
      extra/tcpdump \
      extra/tpm2-tools
  ADD target/ /target/

  RUN mkinitcpio --preset KLoader

FROM scratch as target
  COPY --from=builder /boot/KLoader.efi KLoader.efi
