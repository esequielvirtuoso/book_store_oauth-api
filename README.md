# book_store_oauth-api
OAuth API

This application follows a Domain Driven Development, starting from the Domains (Entities)
to the backward layers (Services, Controllers, and Infrastructure Providers).

Also, it uses the concepts of Dependency Injection. You can read more about it by reading [this book](https://www.packtpub.com/product/hands-on-dependency-injection-in-go/9781789132762).

It also uses the [gin](https://github.com/gin-gonic/gin) HTTP framework to handle the server.

User is the core domain of this application.

![alt text](./doc/images/microservicesDiagram.jpg?raw=true)


## Getting Started

### Prerequisites

- [Golang](http://golang.org/)(>11.0)
- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](http://docker.com)
- [Docker Compose](https://docs.docker.com/compose/install/)


### Environment variables

```
export USER_API_URL=http://localhost:5001 &&

export OAUTH_CASSANDRA_HOST=127.0.0.1 &&

export OAUTH_CASSANDRA_PORT=9043
```


### Installing and running locally
```shell script
# Install dependencies
make install

# Run postgres locally as a container
make env

# Run server locally
make run

# Run server locally with custom environment variables
MYSQL_URL=root:passwd@tcp(127.0.0.1:3305)/users_db?charset=utf8 \
make run
```

### Setting up git hooks

After cloning the repository, change the git hooks path so it's only possible to commit code with the required quality.

```shell script
make git-config
```

## Running the tests and coverage report

To view report of tests locally use the following command:

```bash
make env # prepares environment for testing
make test
```

## Running the lint verification

```bash
make lint
```
_Lint report generate by [GolangCI-lint](https://github.com/golangci/golangci-lint)._

## Running vulnerability check in Go dependencies
```bash
make audit
```
_Audit report generate by [Nancy](https://github.com/sonatype-nexus-community/nancy)._


## Deployment

### Build

```bash
make build
```

### Create release image, add tag and push

Tip: Do not forget to login to docker hub account.
```bash
docker login -u "username" -p "access_token" docker.io
```

```bash
make image tag push
```

### OR

```bash
make image
```
```bash
make tag VERSION=1.0.0
```
```bash
sudo make push VERSION=1.0.0
```

Note: Change the version to the desired version.

### Run registry image locally

Tip: Before running it locally you need to prepare the environment with the following command.
```bash
make env
```
It creates the network and the database containers.

```bash
make run-docker

make remove-docker
```

## Project Structure

### Package organization

The package structure used in this project was inspired by the [golang-standards](https://github.com/golang-standards/project-layout) project.

### Project layers organization

The project layers structure used in this project was inspired by the **Hexagonal Architecture** (Ports & Adapters).


## Contributing
See [CONTRIBUTING](CONTRIBUTING.md) documentation for more details.


## Changelog
See [CHANGELOG](CHANGELOG.md) documentation for more details.
