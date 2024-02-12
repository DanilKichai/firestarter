# syntax=docker/dockerfile:1.2.1

# build target
FROM --platform=linux/amd64 archlinux:base as builder

# install toolchain
RUN pacman \
  --sync \
  --sysupgrade \
  --refresh \
  --noconfirm \
  --needed \
  core/gzip \
  core/linux \
  core/systemd-ukify \
  extra/cpio

# prepare initramfs target directory
WORKDIR /initramfs
RUN \
  mkdir dev && \
  mknod -m 622 dev/console c 5 1 && \
  mknod -m 666 dev/null    c 1 3 && \
  mknod -m 444 dev/random  c 1 8 && \
  mknod -m 444 dev/urandom c 1 9 && \
  mknod -m 666 dev/zero    c 1 5

# install target packages
RUN pacman \
  --root . \
  --dbpath /tmp \
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
  extra/sbsigntools \
  extra/tcpdump \
  extra/tpm2-tools

# install kernel modules
RUN cp --recursive /lib/modules lib/

# install bootstrap
ADD bootstrap/bootstrap.sh opt/bootstrap/
ADD bootstrap/logo opt/bootstrap/
ADD bootstrap/bootstrap@.service lib/systemd/system/

# replace getty with bootstrap service
RUN \
  systemctl --root=. disable getty@tty1.service && \
  systemctl --root=. enable bootstrap@tty1.serivce

# create cpio archive
RUN \
  find | \
  cpio \
    --create \
    --format=newc \
    --owner=root:root \
    --directory=. | \
  gzip \
    --best \
    >/initramfs.cpio.gz

# ukify target
RUN ukify build \
  --linux=/boot/vmlinuz-linux \
  --initrd=/initramfs.cpio.gz \
  --cmdline="rc_init=/lib/systemd/systemd" \
  --os-release="KLoader" \
  --efi-arch="x64" \
  --output=/KLoader.efi

# pick out build target
FROM scratch as target
COPY --from=builder /KLoader.efi KLoader.efi
