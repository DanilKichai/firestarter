.PHONY: all
all: KLoader.efi

.PHONY: builder
builder:
	docker buildx create \
		--name KLoader \
		--node main

logo:
	bootstrap/logo.gen >bootstrap/logo

KLoader.efi: builder logo
	docker buildx build \
		--builder KLoader \
		--progress plain \
		--target target \
		--output "type=local,dest=." \
		.

.PHONY: clean
clean:
	-rm -rf KLoader.efi bootstrap/logo
	-docker buildx rm \
		--force \
		--builder KLoader
