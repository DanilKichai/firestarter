# syntax=docker/dockerfile:experimental

# set defaults
ARG GENTOO_STAGE3_IMAGE="gentoo/stage3:musl-hardened"
ARG GENTOO_PORTAGE_SNAPSHOT=""
ARG WORKSPACE=/build

# build target
FROM "${GENTOO_STAGE3_IMAGE}" as builder

# export arguments to stage environment
ENV GENTOO_PORTAGE_SNAPSHOT="${GENTOO_PORTAGE_SNAPSHOT}"
ENV WORKSPACE="${WORKSPACE}"

# prepare gentoo emerge
RUN \
  mkdir -p /var/db/repos/gentoo && \
    emerge-webrsync $( \
      [ -n "${GENTOO_PORTAGE_SNAPSHOT}" ] && \
        echo -n "--revert=${GENTOO_PORTAGE_SNAPSHOT}" \
    )
ADD package.use /etc/portage/package.use

# build packages
RUN --security=insecure \
  emerge --update --buildpkg --newuse --with-bdeps=y \
    sys-apps/busybox \
    sys-fs/lvm2 \
    sys-libs/zlib \
    sys-fs/cryptsetup

# custom build packages
RUN --security=insecure \
  LDFLAGS="-static" \
    emerge --update --buildpkg --newuse --with-bdeps=y \
      sys-apps/kexec-tools \
      sys-block/sedutil

# install gentoo-source
RUN --security=insecure \
  emerge --update --newuse --with-bdeps=y \
    sys-kernel/gentoo-sources

# prepare initramfs target
WORKDIR "${WORKSPACE}/initramfs"
RUN mkdir \
  dev \
  proc \
  sys \
  boot \
  root
RUN \
  mknod -m 622 dev/console c 5 1 && \
  mknod -m 666 dev/null    c 1 3 && \
  mknod -m 666 dev/zero    c 1 5 && \
  mknod -m 444 dev/random  c 1 8 && \
  mknod -m 444 dev/urandom c 1 9
RUN qtbz2 --tarbz2 --stdout \
  /var/cache/binpkgs/sys-apps/busybox/busybox-*.xpak | \
    tar --extract --zstd
RUN tar --extract --file=usr/share/busybox/busybox-links.tar
RUN qtbz2 --tarbz2 --stdout \
  /var/cache/binpkgs/sys-apps/kexec-tools/kexec-tools-*.xpak | \
    tar --extract --zstd
RUN qtbz2 --tarbz2 --stdout \
  /var/cache/binpkgs/sys-fs/lvm2/lvm2-*.xpak | \
    tar --extract --zstd
RUN \
  for file in sbin/*.static; do \
    mv --force "${file}" "${file::-7}"; \
  done
RUN qtbz2 --tarbz2 --stdout \
  /var/cache/binpkgs/sys-fs/cryptsetup/cryptsetup-*.xpak | \
    tar --extract --zstd
RUN qtbz2 --tarbz2 --stdout \
  /var/cache/binpkgs/sys-block/sedutil/sedutil-*.xpak | \
    tar --extract --zstd
RUN cp /etc/group \
  /etc/passwd \
  etc/
ADD init .

# build kernel
WORKDIR /usr/src/linux
ADD linux.conf .config
RUN \
  echo "CONFIG_BLK_DEV_INITRD=y" >>.config && \
  echo "CONFIG_INITRAMFS_SOURCE=\"${WORKSPACE}/initramfs\"" >>.config && \
  make olddefconfig
RUN make -j "$(cat /proc/cpuinfo | grep processor | wc -l)"

# copy target
FROM scratch
COPY --from=builder /usr/src/linux/arch/x86/boot/bzImage bbloader.efi
