all: integration

DOCKER_IMAGE_TAG :=	$(shell python get_docker_build_version.py)
DOCKER_GIT_RELEASE_TAG := $(shell python get_latest_git_release_tag.py)

build:
	env GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -v -o enchilada server.go

docker: build
	docker build -t www.dockerhub.us/enchilada:$(DOCKER_IMAGE_TAG) .

integration: docker
	./integration/test.sh

publish: integration
	docker tag www.dockerhub.us/enchilada:$(DOCKER_IMAGE_TAG) www.dockerhub.us/enchilada:$(DOCKER_GIT_RELEASE_TAG)
	docker push www.dockerhub.us/enchilada:$(DOCKER_GIT_RELEASE_TAG)

