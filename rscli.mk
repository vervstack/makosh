gen: gen-server-grpc
build-local-container:
	docker buildx build \
			--load \
			--platform linux/arm64 \
			-t makosh:local .

### Grpc server generation
gen-server-grpc: .download-grpc-deps .gen-server-grpc

.download-grpc-deps:
	EASYPPATH=proto_deps easyp mod download

.update-grpc-deps:
	EASYPPATH=proto_deps easyp mod update


.gen-server-grpc:
	EASYPPATH=proto_deps easyp generate