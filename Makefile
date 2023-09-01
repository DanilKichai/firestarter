ARCHLINUX_BASE_IMAGE = "archlinux:base"
LINUX_KERNEL_VERSION = ""

BUILDKIT_STEP_LOG_MAX_SIZE = 104857600

all:          create_builder build_target
olddefconfig: create_builder build_olddefconfig

create_builder:
	docker buildx create \
		--name KLoader \
		--node main \
		--driver-opt "env.BUILDKIT_STEP_LOG_MAX_SIZE=$(BUILDKIT_STEP_LOG_MAX_SIZE)"

build_olddefconfig:
	docker buildx build \
		--builder KLoader \
		--progress plain \
		--build-arg "ARCHLINUX_BASE_IMAGE=$(ARCHLINUX_BASE_IMAGE)" \
		--build-arg "LINUX_KERNEL_VERSION=$(LINUX_KERNEL_VERSION)" \
		--target olddefconfig \
		--output "type=local,dest=." \
		.

build_target:
	docker buildx build \
		--builder KLoader \
		--progress plain \
		--build-arg "ARCHLINUX_BASE_IMAGE=$(ARCHLINUX_BASE_IMAGE)" \
		--build-arg "LINUX_KERNEL_VERSION=$(LINUX_KERNEL_VERSION)" \
		--target target \
		--output "type=local,dest=." \
		.

remove_target:
	rm -rf KLoader.efi

remove_builder:
	docker buildx rm \
		--force \
		--builder KLoader || \
	true

clean: remove_target remove_builder
