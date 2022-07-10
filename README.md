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

You can also just run the cdk tests using

    npm run test:cdk

and for the Golang packages

    npm run test:go

## Deploy
### Deploy to AWS
> Make sure you installed the [AWS Cloud Development Kit](https://aws.amazon.com/cdk/).

    cdk bootstrap
    cdk synth
    cdk deploy

### Deploy as Docker image/container
First, create a docker image using the instruction

    docker build -t login-service .

then you can create and run a container using

    docker run -p 8080:8080 --name login-service-instance -t login-service

The service will be available at `http://localhost:8080`