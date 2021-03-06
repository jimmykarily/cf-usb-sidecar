# Catalog Service Manager

This is the repository for the catalog service manager(aka Side car container).
This repository holds the rest API for Catalog service manager.

<!-- TOC depthFrom:1 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Setting up project environment/dependencies](#setting-up-dependencies)
- [Building](#building)
- [Environment Variables](#environment-variables)
- [Run the service](#run-the-service)
  - [In go](#in-go)
  - [On docker](#run-service-with-examplemysql-extensions)
  - [Helm with SCF](#running-with-the-scf-vagrant-machine)
- [Testing CSM Service](#testing-csm-service)
- [Docker Containers for Catalog Service Manager](#docker-containers-for-catalog-service-manager)

## Setting up Dependencies

Please make sure you have git and mercurial installed on your machine.

### install mercurial on mac

Either you can use brew to install mercurial

```
brew install mercurial
```

or you can also download latest mercurial package from their site.

### install go tools/dependencies for catalog service manager

```
make tools
```

## Building

To build the mysql sidecar and associated helm chart do

```
make build-image
make publish-image

cd csm-extensions/services/dev-mysql
make build-image
make publish-image
make helm
```

Note however that the environment variables `DOCKER_REPOSITORY` and
`DOCKER_ORGANIZATION` have to be set to suitable values before
invoking the above script.

The docker images needed by the sidecar's helm chart will be pushed to
the specified repository and organization. The chart itself will be
stored in the subdirectory `output` of `dev-mysql`.

## Environment Variables

Catalog Service Manager reads it configuration/settings from environment
variables.

### SIDECAR_PARAMETERS

This is the list of environment variables that Catalog Service Manager should
read, this needs to be set if you are running catalog service manager in the
default mode (specifically without create connections extensions). When
extensions are not present, catalog service manager reads this environment and returns
all the variables listed here.

```
export SIDECAR_PARAMETERS="HOSTNAME USERNAME"
```

### SIDECAR_API_KEY

This is a required environment variable that must be set when starting up the
service. The SIDECAR_API_KEY is used to authorize all requests to the service via an
authentication token. If this variable is not set, the service will not start
up.

```
export SIDECAR_API_KEY="sidecar-auth-token"
```

### SIDECAR_HOME

This is the path of the CSM home directory, default value for this is
/catalog-service-manager which is the path on the Catalog Service Manager's base
docker image where all the extensions are available.

```
export SIDECAR_HOME=${GOPATH}/src/github.com/SUSE/cf-usb-sidecar/examples/mysql
```

### SIDECAR_DEBUG

Set this environment variable if you want catalog service manager to generate
logs in Debug mode. By default catalog service manager runs with log level Info,
but for debugging/testing you can set this environment variable and get more
details from catalog service manager logs.

```
export SIDECAR_DEBUG=true
```

### DEV_MODE

Set this environment variable if you want catalog service manager to keep the
output files written by the extensions. By default catalog service
manager runs with DEV_MODE set to off, and it always deletes the output file
where extensions write their output.

```
export DEV_MODE=true
```

### SIDECAR_EXT_TIMEOUT and SIDECAR_EXT_TIMEOUT_ERROR

* SIDECAR_EXT_TIMEOUT has a default value of 30s
* SIDECAR_EXT_TIMEOUT_ERROR has a default value of 2s

Set SIDECAR_EXT_TIMEOUT environment variable to a value that represents the number of seconds that the
catalog service manager will wait for a response from the extension before
sending it a request to stop.

The catalog service manager will wait for a graceful stop a number of
SIDECAR_EXT_TIMEOUT_ERROR seconds. If in this interval the extension did not
stop, the manager will try to force stop the extension.

```
export SIDECAR_EXT_TIMEOUT=30
export SIDECAR_EXT_TIMEOUT_ERROR=2
```

### TLS_CERT_FILE and TLS_PRIVATE_KEY_FILE

CSM will start with TLS automatically and use `TLS_CERT_FILE` and `TLS_PRIVATE_KEY_FILE`
environment variables that are set to the default location of `/etc/secrets/tls-cert-file`
and `/etc/secrets/tls-private-key-file`.

## Run the service

### In go

You can use make run to run the service on command line

```
make run
```

This will start service listening on port https://0.0.0.0:8081 on the machine
where you run this command. You can use above mentioned environment variables
(SIDECAR_HOME and SIDECAR_DEBUG) to alter the behavior of the service. (specifically if
you are debugging an issue or while writing new extensions )


### Run service with example/mysql extensions

1. To start with you need mysql instance to which you can connect to, if you
have one already you can skip this first step, if you don't then run following
commands

```
cd  example/mysql
make tools
```

This should start a mysql container and bind port 3306 on the docker host to
mysql port on the container After this is done, you should be able to connect to
the mysql with

```
MY_DOCKER_HOST_IP=`env | grep DOCKER_HOST | cut -d "/" -f 3 | cut -d ":" -f 1`
mysql -h ${MY_DOCKER_HOST_IP} -P 3306 -u root -proot123

```

2. Set following environment variables

```
export SIDECAR_HOME=`pwd`/examples/mysql/SIDECAR_HOME
export MYSQL_SERVICE_HOST=`env | grep DOCKER_HOST | cut -d "/" -f 3 | cut -d ":" -f 1`
export MYSQL_SERVICE_PORT_MYSQL=3306
export MYSQL_ROOT_PASSWORD=root123
```

*optional override of token when testing service*


```
export SIDECAR_API_KEY=my-foo-bar-key
```


3. Run the make command to start the service
```
make run
```

### Running with the SCF Vagrant machine

For a simplified setup with the SCF vagrant machine (or any other compatible
Kubernetes-based deployment), use `scripts/cf-deploy.sh`:

```bash
scripts/cf-deploy.sh mysql
# Alternatively, `scripts/cf-deploy.sh postgres`
```

This will deploy the helm chart (which should be manually built before hand)
using the current configuration.

## Testing CSM Service

If authorization was enabled with a CMS_API_KEY environment variable specified,
then supply the token in a header called `x-sidecar-token` on each call.

### create workspace

```bash
curl -k -X POST http://localhost:8081/workspaces -H "x-sidecar-token: sidecar-auth-token" -H "content-type: application/json" -d '{"workspace_id":"test_workspace"}'

```

### get workspace

```bash
curl -k -X GET http://localhost:8081/workspaces/test_workspace -H "x-sidecar-token: sidecar-auth-token"
```

### create credentials

```bash
curl -k -X POST http://localhost:8081/workspaces/test_workspace/connections -H "x-sidecar-token: sidecar-auth-token" -H "content-type: application/json" -d '{"connection_id":"test_user"}'
```

### get credentials

```bash
curl -k -X GET http://localhost:8081/workspaces/test_workspace/connections/test_user -H "x-sidecar-token: sidecar-auth-token"
```

### delete credentials

```bash
curl -k -X DELETE http://localhost:8081/workspaces/test_workspace/connections/test_user -H "x-sidecar-token: sidecar-auth-token"
```

### delete workspace

```bash
curl -k -X DELETE http://localhost:8081/workspaces/test_workspace -H "x-sidecar-token: sidecar-auth-token"
```

## Docker Containers for Catalog Service Manager

CSM provides a base docker image which users creating CSM for their services
should be able to extended easily. For more information on how to extend the
base image please refer to the documentation on Extending_CSM.md

### catalog-service-manager:build image (based on ubuntu:14.04)

This is the build image which has all the CSM dependencies installed (mercurial,
swagger, make, git etc) on it. By having this image locally available we save at
least few minutes, every time we have to build other images in the workflow.
Although you don't need to build this image explicitly its built by the other
scripts if they don't find it in docker images.

```
scripts/generate-build-image.sh
```

### catalog-service-manager:development (based on catalog-service-manager:build)

This is the dev image for CSM and primarily used by developers working on CSM.
It uses catalog-service-manager:build and copies CSM code and compiles it there
and runs tests on the container. If required this container can be started as a
docker container locally with the command

```
make dev-base
```

or if you want to run the script directly

```
scripts/development/generate-development-base-image.sh
```

This also publishes the image  with tag catalog-service-manager:base

### catalog-service-manager:release (based on catalog-service-manager:build)

This is the intermediate docker image for CSM and primarily used by for
generating the release artifacts. It uses catalog-service-manager:build and
copies CSM code and compiles it there and runs tests on the container. When the
image building is finished we start the container with a local directory mounted
on the container where we copy the build artifacts

```
make release-base
```

or if you want to run the script directly

```
scripts/release/generate-release-base-image.sh
```

This also publishes the image with tag catalog-service-manager:base, and this is
the image/tag which will be used by all the services.
