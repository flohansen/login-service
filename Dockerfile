FROM alpine:3.16.0

ARG GOLANG_VERSION=1.18.2

ENV LOGIN_SERVICE_HOST=0.0.0.0
ENV LOGIN_SERVICE_PORT=8080
ENV LOGIN_SERVICE_JWT_SIGN_KEY=secret
ENV LOGIN_SERVICE_DATABASE_HOST=127.0.0.1
ENV LOGIN_SERVICE_DATABASE_PORT=5432
ENV LOGIN_SERVICE_DATABASE_USER=username
ENV LOGIN_SERVICE_DATABASE_PASSWORD=password
ENV LOGIN_SERVICE_DATABASE_NAME=database

# Install required packages
RUN apk update
RUN apk add go gcc bash musl-dev openssl-dev ca-certificates
RUN update-ca-certificates

# Install go
RUN wget https://dl.google.com/go/go${GOLANG_VERSION}.src.tar.gz
RUN tar -C /usr/local -xzf go${GOLANG_VERSION}.src.tar.gz
RUN cd /usr/local/go/src && ./make.bash
ENV PATH=$PATH:/usr/local/go/bin
RUN rm -f go${GOLANG_VERSION}.src.tar.gz
RUN apk del go

# Build the project
COPY . .
RUN go build -o build/app ./src/main.go

EXPOSE ${LOGIN_SERVICE_PORT}
CMD [ "build/app" ]