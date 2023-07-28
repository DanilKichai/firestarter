GENTOO_BASE_IMAGE = "gentoo/stage3:musl-hardened"
BUILDKIT_STEP_LOG_MAX_SIZE = 104857600

all: create_builder build_bbloader

create_builder:
	docker buildx create \
		--name bbloader-builder \
		--node bbloader-builder \
		--buildkitd-flags "--allow-insecure-entitlement security.insecure" \
		--driver-opt "env.BUILDKIT_STEP_LOG_MAX_SIZE=${BUILDKIT_STEP_LOG_MAX_SIZE}"

build_bbloader:
	docker buildx build \
		--builder bbloader-builder \
		--allow security.insecure \
		--progress plain \
		--build-arg "GENTOO_BASE_IMAGE=$(GENTOO_BASE_IMAGE)" \
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
