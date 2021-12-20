# Installing

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone git clone https://github.com/nakiner/faceit.git

#move to project
$ cd faceit

# Build the docker image first (or skip, to use docker.io prebuilt one)
$ make docker

# Run the application
$ make docker-run

# check if the containers are running
$ docker ps

# Execute the call
$ curl localhost:8080/liveness

# Run tests (optional)
$ make test
```

# Explanation

Based on my development experience with Go I have decided to use go-kit as main toolkit for maintaining all access 
and reaching observability to every component in system. 

Service consists of multiple layers:

Based on https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

- Models Layer
- Repository Layer (internal/repository)
- Service Layer (business logic pkg/user/service.go)
- Transport layer (HTTP/GRPC/MQ)

Own codegen solution used to generate boilerplate code to get exactly: 

- HTTP listener
- GRPC listener
- HTTP client
- GRPC client
- Service layer boilerplate
- Basic observability

For each dependency such as SQL/MQ written unit-tests and in addition in (test/integration) written tests to perform integration checks for main deps;

As for personal preference I think current layer set satisfy service requirements and logic disturbed over all layers to bring scalability and improve 
understanding of code in this project.

In addition, common project structure https://github.com/golang-standards/project-layout used to store all code in right places;

# Possible addons

- There is no perfect service with all satisfied requirements so in current instance we can introduce more service layers (ex. for different user actions)  create 
new division to serve logic and handle properties properly;
- For current service would be a good call to reform repository contents and add more connectors to serve more data sources (ex. MongoDB) and disturb logic over connectors;
- Write more tests: not all parts tested, would be great to test all layers and see failures before they come into production. For now integrity testing solves this problem but requires starting instances to perform testing;
- Add manifests to perform orchestration (k8s) to verify integrity upon release;