# include help.mk

# tell Make the following targets are not files, but targets within Makefile
.PHONY: version clean audit lint install build image tag push release run run-local remove-docker env env-stop print-var check-env check-used-ports
.DEFAULT_GOAL := help

GITHUB_GROUP = esequielvirtuoso
HUB_HOST     = hub.docker.com/repository/docker
HUB_USER 	 = esequielvirtuoso
HUB_REPO     = book_store_oauth-api

BUILD         	= $(shell git rev-parse --short HEAD)
DATE          	= $(shell date -uIseconds)
VERSION  	  	= $(shell git describe --always --tags)
NAME           	= $(shell basename $(CURDIR))
IMAGE          	= $(HUB_USER)/$(HUB_REPO):$(BUILD)

CASSANDRA_NAME = cassandra_$(NAME)_$(BUILD)

# NETWORK_NAME can be dinamically generated with the following env set
# NETWORK_NAME  = network_$(NAME)_$(BUILD)
# However, we have set it up with a static name to simplify the local
# connection tests between the apps containers
NETWORK_NAME = network_book_store

USER_API_URL=http://book_store_users-api:8081
OAUTH_CASSANDRA_HOST=cassandra
OAUTH_CASSANDRA_PORT="9042"

check-used-ports:
	sudo netstat -tulpn | grep LISTEN

print_var:
	echo $(DATE)

git-config:
	git config --replace-all core.hooksPath .githooks

check-env-%:
	@ if [ "${${*}}" = ""  ]; then \
		echo "Variable '$*' not set"; \
		exit 1; \
	fi

version: ##@other Check version.
	@echo $(VERSION)

clean: ##@dev Remove folder vendor, public and coverage.
	rm -rf vendor public coverage

install: clean ##@dev Download dependencies via go mod.
	GO111MODULE=on go mod download
	GO111MODULE=on go mod vendor

audit: ##@check Run vulnerability check in Go dependencies.
	DOCKER_BUILDKIT=1 docker build --progress=plain --target=audit --file=Dockerfile .

lint: ##@check Run lint on docker.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain \
		--target=lint \
		--file=Dockerfile .

env: ##@environment Create network and run CASSANDRA container.
	CASSANDRA_NAME=${CASSANDRA_NAME} \
	NETWORK_NAME=${NETWORK_NAME} \
	docker-compose up -d

env-ip: ##@environment Return local CASSANDRA IP (from Docker container)
	@echo $$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${CASSANDRA_NAME})

env-stop: ##@environment Remove CASSANDRA container and remove network.
	CASSANDRA_NAME=${CASSANDRA_NAME} NETWORK_NAME=${NETWORK_NAME} docker-compose kill
	CASSANDRA_NAME=${CASSANDRA_NAME} NETWORK_NAME=${NETWORK_NAME} docker-compose rm -vf
	docker network rm $(NETWORK_NAME)

test: ##@check Run tests and coverage.
	docker build --progress=plain \
		--network $(NETWORK_NAME) \
		--tag $(IMAGE) \
		--build-arg USER_API_URL="${USER_API_URL}" \
		--build-arg OAUTH_CASSANDRA_HOST="${OAUTH_CASSANDRA_HOST}" \
		--build-arg OAUTH_CASSANDRA_PORT=${OAUTH_CASSANDRA_PORT} \
		--target=test \
		--file=Dockerfile .

	-mkdir coverage
	docker create --name $(NAME)-$(BUILD) $(IMAGE)
	docker cp $(NAME)-$(BUILD):/index.html ./coverage/.
	docker rm -vf $(NAME)-$(BUILD)

build: ##@build Build image.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain \
		--tag $(IMAGE) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		--target=build \
		--file=Dockerfile .

image: check-env-VERSION ##@build Create release docker image.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain \
		--tag $(IMAGE) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		--target=image \
		--file=Dockerfile .

tag: check-env-VERSION ##@build Add docker tag.
	docker tag $(IMAGE) \
		$(HUB_USER)/$(HUB_REPO):$(VERSION)

push: check-env-VERSION ##@build Push docker image to registry.
	docker push $(HUB_USER)/$(HUB_REPO):$(VERSION)

release: check-env-TAG ##@build Create and push git tag.
	git tag -a $(TAG) -m "Generated release "$(TAG)
	git push origin $(TAG)

run:
	go run src/main.go

run-local: ##@dev Run locally.
	USER_API_URL="${USER_API_URL}" \
	OAUTH_CASSANDRA_HOST="${OAUTH_CASSANDRA_HOST}" \
	OAUTH_CASSANDRA_PORT=${OAUTH_CASSANDRA_PORT} \
	run

run-docker: check-env-USER_API_URL check-env-OAUTH_CASSANDRA_HOST check-env-OAUTH_CASSANDRA_PORT ##@docker Run docker container.
	docker run --rm \
		--name $(NAME) \
		--network $(NETWORK_NAME) \
		-e LOGGER_LEVEL=debug \
		-e USER_API_URL="${USER_API_URL}" \
		-e OAUTH_CASSANDRA_HOST="${OAUTH_CASSANDRA_HOST}" \
		-e OAUTH_CASSANDRA_PORT=${OAUTH_CASSANDRA_PORT} \
		-p 5002:8082 \
		$(IMAGE)

remove-docker: ##@docker Remove docker container.
	-docker rm -vf $(NAME)
