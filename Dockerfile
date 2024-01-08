# syntax=docker/dockerfile:1.2.1

# set defaults
ARG ARCHLINUX_BASE_IMAGE="archlinux:base"
ARG ARCHLINUX_TOOLCHAIN_PACKAGES=" \
  core/base-devel \
  core/gettext \
  core/libelf \
  core/links \
  core/perl \
  core/python \
  core/sudo \
  core/tar \
  core/xz \
  extra/bc \
  extra/cpio \
  extra/git \
  extra/pahole \
  extra/wget \
"
ARG LINUX_KERNEL_VERSION=""
ARG INITRMFS_TARGET_PACKAGES=" \
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
  extra/tpm2-tools \
  extra/haveged \
"
ARG INITRMFS_TARGET_AUR_PACKAGES=" \
  sedutil \
"

# build toolchain
FROM "${ARCHLINUX_BASE_IMAGE}" as toolchain
ARG ARCHLINUX_TOOLCHAIN_PACKAGES
RUN pacman \
  --sync \
  --sysupgrade \
  --refresh \
  --noconfirm \
  --needed \
  ${ARCHLINUX_TOOLCHAIN_PACKAGES}
RUN \
  useradd makepkg --create-home && \
  echo 'makepkg ALL=(ALL) NOPASSWD: ALL' >>/etc/sudoers.d/makepkg
RUN \
  mkdir /tmp/yay && \
  cd /tmp/yay && \
  git clone https://aur.archlinux.org/yay.git . && \
  chown -R makepkg * . && \
  sudo -u makepkg makepkg \
    --syncdeps \
    --noconfirm \
    --install

# build pre-kernel
FROM toolchain as pre-kernel
WORKDIR /usr/src
ARG LINUX_KERNEL_VERSION
RUN \
  VERSION="$( \
    if [ -z "${LINUX_KERNEL_VERSION}" ]; then \
      links -dump https://www.kernel.org/ | \
        grep 'longterm:' | \
          awk '{ print $2; exit }'; \
    else \
      echo "${LINUX_KERNEL_VERSION}"; \
    fi \
  )" && \
  MAJOR="$( \
    echo "${VERSION}" | \
      awk -F '.' '{ print $1 }' \
  )" && \
  wget "https://cdn.kernel.org/pub/linux/kernel/v${MAJOR}.x/linux-${VERSION}.tar.xz" && \
  tar --extract --file "linux-${VERSION}.tar.xz" && \
  chown --recursive "$(id -u):$(id -g)" "linux-${VERSION}" && \
  ln --symbolic --force "linux-${VERSION}" linux
WORKDIR /usr/src/linux
RUN make mrproper
ADD linux.conf .config
RUN \
  echo 'CONFIG_BLK_DEV_INITRD=y' >>.config && \
  echo 'CONFIG_INITRAMFS_SOURCE="/initramfs"' >>.config && \
  echo 'CONFIG_MODULES=y' >>.config && \
  make olddefconfig

# build initramfs
FROM pre-kernel as initramfs
WORKDIR /initramfs
RUN \
  mkdir dev && \
  mknod -m 622 dev/console c 5 1 && \
  mknod -m 666 dev/null    c 1 3 && \
  mknod -m 444 dev/random  c 1 8 && \
  mknod -m 444 dev/urandom c 1 9 && \
  mknod -m 666 dev/zero    c 1 5
RUN mkdir /tmp/pacman
ARG INITRMFS_TARGET_PACKAGES
RUN \
  [[ "${INITRMFS_TARGET_PACKAGES}" =~ ^[" "]*$ ]] || \
    pacman \
      --root /initramfs \
      --dbpath /tmp/pacman \
      --sync \
      --sysupgrade \
      --refresh \
      --noconfirm \
      --needed \
      ${INITRMFS_TARGET_PACKAGES}
ARG INITRMFS_TARGET_AUR_PACKAGES
RUN \
  [[ "${INITRMFS_TARGET_AUR_PACKAGES}" =~ ^[" "]*$ ]] || \
    sudo -u makepkg yay \
      --root /initramfs \
      --dbpath /tmp/pacman \
      --sync \
      --noconfirm \
      ${INITRMFS_TARGET_AUR_PACKAGES}
RUN \
  ln --symbolic --force /lib/systemd/systemd init && \
  mv \
    etc/systemd/system/getty.target.wants/getty@tty1.service \
    etc/systemd/system/getty.target.wants/getty@console.service && \
  chroot . systemctl enable haveged.service
ADD getty@console.conf etc/systemd/system/getty@console.service.d/override.conf
ADD bootstrap usr/local/bin/
ADD issue etc/

# build kernel
FROM initramfs as kernel
WORKDIR /usr/src/linux
RUN \
  make --jobs=$(nproc) && \
  make --jobs=$(nproc) INSTALL_MOD_PATH=/initramfs modules_install && \
  make --jobs=$(nproc) bzImage

# pick out kernel olddefconfig
FROM scratch as olddefconfig
COPY --from=pre-kernel /usr/src/linux/.config linux.conf.grand

# pick out build target
FROM scratch as target
COPY --from=kernel /usr/src/linux/arch/x86/boot/bzImage KLoader.efi
