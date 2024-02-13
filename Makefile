.PHONY: all
all: KLoader.efi

.PHONY: builder
builder:
	docker buildx create \
		--name KLoader \
		--node main

KLoader.efi: builder
	docker buildx build \
		--builder KLoader \
		--progress plain \
		--target target \
		--output "type=local,dest=." \
		.

.PHONY: clean
clean:
	-rm -rf KLoader.efi
	-docker buildx rm \
		--force \
		--builder KLoader
