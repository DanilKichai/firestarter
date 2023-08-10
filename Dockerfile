# syntax=docker/dockerfile:experimental

# set defaults
ARG GENTOO_STAGE3_IMAGE="gentoo/stage3:musl"
ARG GENTOO_PORTAGE_SNAPSHOT=""
ARG PORTAGE_UNWANTED_PACKAGES=" \
  virtual/udev \
"
ARG TARGET_ADDITION_PACKAGES=" \
  app-crypt/sbsigntools \
  app-crypt/tpm2-tools \
  sys-apps/kexec-tools \
  sys-block/sedutil \
  sys-boot/efibootmgr \
  sys-fs/cryptsetup \
  sys-fs/eudev \
  sys-fs/lvm2 \
"

# build portage
FROM "${GENTOO_STAGE3_IMAGE}" as portage
ARG GENTOO_PORTAGE_SNAPSHOT
ARG PORTAGE_UNWANTED_PACKAGES
ENV GENTOO_PORTAGE_SNAPSHOT="${GENTOO_PORTAGE_SNAPSHOT}"
ENV PORTAGE_UNWANTED_PACKAGES="${PORTAGE_UNWANTED_PACKAGES}"
RUN \
  mkdir -p /var/db/repos/gentoo && \
    emerge-webrsync $( \
      [ -n "${GENTOO_PORTAGE_SNAPSHOT}" ] && \
        echo -n "--revert=${GENTOO_PORTAGE_SNAPSHOT}" \
    )
ADD package.use /etc/portage/package.use
RUN emerge --unmerge ${PORTAGE_UNWANTED_PACKAGES}

# build packages
FROM portage as packages
ARG TARGET_ADDITION_PACKAGES
ENV TARGET_ADDITION_PACKAGES="${TARGET_ADDITION_PACKAGES}"
RUN --security=insecure \
  emerge --buildpkg --emptytree \
    sys-apps/busybox \
    sys-kernel/gentoo-sources \
    ${TARGET_ADDITION_PACKAGES}

# build initramfs
FROM portage as initramfs
ARG TARGET_ADDITION_PACKAGES
ENV TARGET_ADDITION_PACKAGES="${TARGET_ADDITION_PACKAGES}"
COPY --from=packages /var/cache/binpkgs /var/cache/binpkgs
WORKDIR /initramfs
RUN mkdir \
  dev \
  etc \
  proc \
  root \
  sys
RUN \
  mknod -m 622 dev/console c 5 1 && \
  mknod -m 666 dev/null    c 1 3 && \
  mknod -m 444 dev/random  c 1 8 && \
  mknod -m 444 dev/urandom c 1 9 && \
  mknod -m 666 dev/zero    c 1 5
RUN cp /etc/group /etc/passwd etc/
RUN ROOT=/initramfs emerge --getbinpkgonly ${TARGET_ADDITION_PACKAGES}
ADD init .

# build kernel
FROM portage as kernel
COPY --from=initramfs /initramfs         /initramfs
COPY --from=packages  /var/cache/binpkgs /var/cache/binpkgs
RUN emerge --getbinpkgonly sys-kernel/gentoo-sources
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
