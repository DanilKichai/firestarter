# syntax=docker/dockerfile:experimental

# set defaults
ARG GENTOO_STAGE3_IMAGE="gentoo/stage3:musl"
ARG GENTOO_PORTAGE_SNAPSHOT=""
ARG EMERGE_DEFAULT_OPTS=" \
  --accept-properties=-interactive \
  --jobs \
  --newuse \
  --deep \
  --update \
  --oneshot \
"
ARG TARGET_CORE_PACKAGES=" \
  sys-apps/busybox \
  sys-apps/kexec-tools \
  sys-fs/eudev \
"
ARG TARGET_ADDITION_PACKAGES=" \
  app-crypt/sbsigntools \
  app-crypt/tpm2-tools \
  sys-block/sedutil \
  sys-boot/efibootmgr \
  sys-fs/btrfs-progs \
  sys-fs/cryptsetup \
  sys-fs/dosfstools \
  sys-fs/e2fsprogs \
  sys-fs/lvm2 \
"

# build portage
FROM "${GENTOO_STAGE3_IMAGE}" as portage
ARG GENTOO_PORTAGE_SNAPSHOT
ENV GENTOO_PORTAGE_SNAPSHOT="${GENTOO_PORTAGE_SNAPSHOT}"
ARG EMERGE_DEFAULT_OPTS
ENV EMERGE_DEFAULT_OPTS="${EMERGE_DEFAULT_OPTS}"
RUN \
  mkdir -p /var/db/repos/gentoo && \
  emerge-webrsync $( \
    [ -n "${GENTOO_PORTAGE_SNAPSHOT}" ] && \
      echo -n "--revert=${GENTOO_PORTAGE_SNAPSHOT}" \
  )
ADD package.use /etc/portage/package.use

# build initramfs
FROM portage as initramfs
WORKDIR /initramfs
ARG TARGET_CORE_PACKAGES
ENV TARGET_CORE_PACKAGES="${TARGET_CORE_PACKAGES}"
ARG TARGET_ADDITION_PACKAGES
ENV TARGET_ADDITION_PACKAGES="${TARGET_ADDITION_PACKAGES}"
RUN --security=insecure \
  emerge \
      ${TARGET_CORE_PACKAGES} \
      ${TARGET_ADDITION_PACKAGES}
RUN --security=insecure \
  emerge \
    --root="/initramfs" \
    --quickpkg-direct=y \
    --quickpkg-direct-root="/" \
    --with-bdeps=n \
    --implicit-system-deps=n \
      ${TARGET_CORE_PACKAGES} \
      ${TARGET_ADDITION_PACKAGES}
RUN mkdir \
  dev \
  proc \
  root \
  sys
RUN \
  mknod -m 622 dev/console c 5 1 && \
  mknod -m 666 dev/null    c 1 3 && \
  mknod -m 444 dev/random  c 1 8 && \
  mknod -m 444 dev/urandom c 1 9 && \
  mknod -m 666 dev/zero    c 1 5
ADD init .

# build kernel
FROM initramfs as kernel
RUN emerge sys-kernel/gentoo-sources
WORKDIR /usr/src/linux
ADD linux.conf .config
RUN \
  echo 'CONFIG_BLK_DEV_INITRD=y' >>.config && \
  echo 'CONFIG_INITRAMFS_SOURCE="/initramfs"' >>.config && \
  make olddefconfig
RUN make -j "$(cat /proc/cpuinfo | grep processor | wc -l)"

# pick out target
FROM scratch
COPY --from=kernel /usr/src/linux/arch/x86/boot/bzImage bbloader.efi
