# syntax=docker/dockerfile:1.2.1

FROM --platform=linux/amd64 golang:latest as builder
  COPY src/bootstrap /usr/src/bootstrap
  RUN \
    CGO_ENABLED="0" \
      go build \
        -C /usr/src/bootstrap \
        -ldflags='-extldflags=-static' \
        -o /usr/local/bin/bootstrap \
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
      core/lvm2 \
      core/nano \
      core/sed \
      core/traceroute \
      core/systemd \
      core/udev \
      core/util-linux \
      extra/kexec-tools \
      extra/qrencode \
      extra/tcpdump
  ADD target/ /target/

  COPY --from=builder /usr/local/bin/bootstrap /target/usr/local/bin/
  RUN echo -e "debug\ndebug" | passwd --root /target

  RUN mkinitcpio --preset wrapper

FROM scratch as target
  COPY --from=wrapper /boot/wrapper.efi firestarter.efi
