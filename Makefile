all: integration

DOCKER_IMAGE_TAG :=	$(shell python get_docker_build_version.py)
DOCKER_GIT_RELEASE_TAG := $(shell python get_latest_git_release_tag.py)
DOCKER_REPO=www.dockerhub.us
DOCKER_IMAGE=testify

build:
	env GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -v -o $(DOCKER_IMAGE) testify.go

docker: build
	docker build -t $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG) .

integration: docker
	./integration/test.sh

do_test: integration
	python ./integration/test.py

publish: integration
	docker tag $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG) $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_GIT_RELEASE_TAG)
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_GIT_RELEASE_TAG)

