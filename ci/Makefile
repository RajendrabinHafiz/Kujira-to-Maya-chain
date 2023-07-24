
all: builder runner

builder:
	docker buildx build -t registry.gitlab.com/mayachain/mayanode:builder-v4 -f Dockerfile.builder .

runner:
	docker buildx build -t registry.gitlab.com/mayachain/mayanode:runner-base-v1 -f Dockerfile.runner .

push-builder:
	docker push registry.gitlab.com/mayachain/mayanode:builder-v4

push-runner:
	docker push registry.gitlab.com/mayachain/mayanode:runner-base-v1
