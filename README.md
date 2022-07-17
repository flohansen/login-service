# Fitter Login Service

## Setup
### Setup the CDK project
To setup the CDK use the following command set.

    cd cdk
    npm install

### Setup the Golang project
To setup the Golang project use the following command set.

    go get ./src/...

## Test
> Make sure you installed and run [Docker](https://docker.com/). It is used by
> integration tests to start and stop containers automatically.

To run all of the tests (unit, integration), you can run the `test` script.

    npm run test

and for the Golang packages

    npm run test:go

## Deploy

### Deploy as Docker image/container
First, create a docker image using the instruction

    docker build -t login-service .

then you can create and run a container using

    docker run -dp 8080:8080 -e LOGIN_SERVICE_PORT=8080 --name login-service-instance -t login-service

The service will be available at `http://localhost:8080`. Additionally, you may
want to set some environment variables, that change the behaviour of the
service. You should change at least the default password for the database
connection due to security reasons.