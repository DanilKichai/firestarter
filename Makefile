.PHONY: all
all: archshell.efi

.PHONY: builder
builder:
	docker buildx create \
		--name archshell \
		--node main

archshell.efi: builder
	docker buildx build \
		--builder archshell \
		--progress plain \
		--file build/package/Dockerfile \
		--target target \
		--output "type=local,dest=." \
		.

.PHONY: clean
clean:
	-rm -rf archshell.efi
	-docker buildx rm \
		--force \
		--builder archshell
