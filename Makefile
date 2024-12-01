.PHONY: all
all: firestarter.efi

.PHONY: builder
builder:
	docker buildx create \
		--name firestarter \
		--node main

firestarter.efi: builder
	docker buildx build \
		--builder firestarter \
		--progress plain \
		--file build/package/Dockerfile \
		--target target \
		--output "type=local,dest=." \
		.

.PHONY: clean
clean:
	-rm -rf firestarter.efi
	-docker buildx rm \
		--force \
		--builder firestarter
