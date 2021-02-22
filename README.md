
# Integration tests

Various integration tests that can be used on common platforms  


### Release
Stable versions are released via tags that published the docker images to github via actions.

# RabbitMQ integration
RabbitMQ integration is a small test suite to test rabbitmq server installations.  
It works by connecting to a rabbitmq server via a provided AMQP string, it will then create a queue , post a message and consume the message.

## Dockerfiles
Download the dockerfiles where VERSION is any of the release tags. 
```bash
docker pull ghcr.io/barkardk/rabbitmq-client:VERSION
```

The client needs a running rabbitmq server to start up properly   
 


## Installation

Build a test binary , compile a docker image and push to docker registry
```bash
make build
```
Deploy to kubernetes
```bash
kubectl apply -f it/testdata
```
## Usage
Run locally using docker compose
```bash
docker-compose up
```

When deploying to kubernetes the rabbit mq client pod will run as a job, check the job logs for output
```bash
kubectl logs -l app=rabbitmq-client
```

## Parameters
|   Parameter | Default   |  
|:---|---|
| RABBITMQ_AMQP_CONN_STR  | amqp://guest:guest@localhost:5672/  |  
| VERSION  |  git describe --tags --dirty --match='v*' 2>/dev/null || echo v0.0.0) | cut -c2- |  
| DOCKER_REGISTRY | ghcr.io/barkardk  |


![Octocat](https://github.githubassets.com/images/icons/emoji/octocat.png)