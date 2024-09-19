# syntax=docker/dockerfile:1.2.1

FROM --platform=linux/amd64 golang:latest as builder
  COPY src/bootstrap /usr/src/bootstrap
  RUN \
    CGO_ENABLED="0" \
      go build \
        -C /usr/src/bootstrap \
        -ldflags='-extldflags=-static' \
        -o /opt/firestarter/bootstrap \
        .

FROM --platform=linux/amd64 archlinux:base as wrapper
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

  COPY --from=builder /opt/firestarter/bootstrap /target/opt/firestarter/

  RUN mkinitcpio --preset wrapper

FROM scratch as target
  COPY --from=wrapper /boot/wrapper.efi firestarter.efi
