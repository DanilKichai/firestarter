GENTOO_PORTAGE_SNAPSHOT = 20230821
GENTOO_STAGE3_IMAGE = "gentoo/stage3:hardened-20230821"

BUILDKIT_STEP_LOG_MAX_SIZE = 104857600

all: create_builder build_bbloader

create_builder:
	docker buildx create \
		--name bbloader-builder \
		--node bbloader-builder \
		--driver-opt "env.BUILDKIT_STEP_LOG_MAX_SIZE=$(BUILDKIT_STEP_LOG_MAX_SIZE)"

build_bbloader:
	docker buildx build \
		--builder bbloader-builder \
		--progress plain \
		--build-arg "GENTOO_STAGE3_IMAGE=$(GENTOO_STAGE3_IMAGE)" \
		--build-arg "GENTOO_PORTAGE_SNAPSHOT=$(GENTOO_PORTAGE_SNAPSHOT)" \
		--output "type=local,dest=." \
		.

remove_bbloader:
	rm -rf bbloader.efi

remove_builder:
	docker buildx rm \
		--force \
		--builder bbloader-builder || \
	true

clean: remove_bbloader remove_builder
