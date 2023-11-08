ARCHLINUX_BASE_IMAGE = "archlinux:base"
LINUX_KERNEL_VERSION = ""

BUILDKIT_STEP_LOG_MAX_SIZE = 104857600

.PHONY: all
all: KLoader.efi

.PHONY: builder
builder:
	docker buildx create \
		--name KLoader \
		--node main \
		--driver-opt "env.BUILDKIT_STEP_LOG_MAX_SIZE=$(BUILDKIT_STEP_LOG_MAX_SIZE)"

.PHONY: olddefconfig
olddefconfig: builder
	docker buildx build \
		--builder KLoader \
		--progress plain \
		--build-arg "ARCHLINUX_BASE_IMAGE=$(ARCHLINUX_BASE_IMAGE)" \
		--build-arg "LINUX_KERNEL_VERSION=$(LINUX_KERNEL_VERSION)" \
		--target olddefconfig \
		--output "type=local,dest=." \
		.

KLoader.efi: builder
	docker buildx build \
		--builder KLoader \
		--progress plain \
		--build-arg "ARCHLINUX_BASE_IMAGE=$(ARCHLINUX_BASE_IMAGE)" \
		--build-arg "LINUX_KERNEL_VERSION=$(LINUX_KERNEL_VERSION)" \
		--target target \
		--output "type=local,dest=." \
		.

.PHONY: clean
clean:
	-rm -rf KLoader.efi
	-docker buildx rm \
		--force \
		--builder KLoader
