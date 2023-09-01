ARCHLINUX_BASE_IMAGE = "archlinux:base"
LINUX_KERNEL_VERSION = ""

BUILDKIT_STEP_LOG_MAX_SIZE = 104857600

all:          create_builder build_kloader
olddefconfig: create_builder build_olddefconfig

create_builder:
	docker buildx create \
		--name kloader-builder \
		--node kloader-builder \
		--driver-opt "env.BUILDKIT_STEP_LOG_MAX_SIZE=$(BUILDKIT_STEP_LOG_MAX_SIZE)"

build_olddefconfig:
	docker buildx build \
		--builder kloader-builder \
		--progress plain \
		--build-arg "ARCHLINUX_BASE_IMAGE=$(ARCHLINUX_BASE_IMAGE)" \
		--build-arg "LINUX_KERNEL_VERSION=$(LINUX_KERNEL_VERSION)" \
		--target olddefconfig \
		--output "type=local,dest=." \
		.

build_kloader:
	docker buildx build \
		--builder kloader-builder \
		--progress plain \
		--build-arg "ARCHLINUX_BASE_IMAGE=$(ARCHLINUX_BASE_IMAGE)" \
		--build-arg "LINUX_KERNEL_VERSION=$(LINUX_KERNEL_VERSION)" \
		--target kloader \
		--output "type=local,dest=." \
		.

remove_kloader:
	rm -rf kloader.efi

remove_builder:
	docker buildx rm \
		--force \
		--builder kloader-builder || \
	true

clean: remove_kloader remove_builder
