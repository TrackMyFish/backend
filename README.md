# Backend

Provides a gRPC and optional HTTP backend for TrackMyFish.


# Pre-requisites

- Install [grpcurl](https://github.com/fullstorydev/grpcurl)
- Go >= 1.16. Because we use embed, any versions older than 1.16 won't work.

# gRPC requests

## List services

```
grpcurl -plaintext localhost:8080 list
```

## List RPC endpoints

Note: This assumes the api repository is cloned at the same location as this repository (`../api`), that the required dependencies have been install into the `.cache` directory, of the api repository, and that you're using an Apple device (Darwin).

```
grpcurl -protoset <(cd ../proto; ../proto/.cache/Darwin/x86_64/bin/buf image build -o -) -plaintext localhost:8080 list trackmyfish.v1alpha1.TrackMyFishService
```


# HTTP requests

## Add Fish

```
curl -H "Content-Type: application/json" -X POST localhost:8443/api/v1alpha1/fish -d '{"genus": "Pterophyllum", "species": "scalare", "commonName": "Angel Fish", "gender": "MALE"}'
```

## List Fish

```
curl -H "Content-Type: application/json" -X GET localhost:8443/api/v1alpha1/fish
```

## Delete Fish

```
curl -H "Content-Type: application/json" -X DELETE localhost:8443/api/v1alpha1/fish/1
```

## Add Tank Statistic

```
curl -H "Content-Type: application/json" -X POST localhost:8443/api/v1alpha1/tank/statistics -d '{"testDate": "2021/08/06 10:00", "ammonia": "2.0"}'
```

## List Tank Statistics

```
curl -H "Content-Type: application/json" -X GET localhost:8443/api/v1alpha1/tank/statistics
```

## Delete Tank Statistics

```
curl -H "Content-Type: application/json" -X DELETE localhost:8443/api/v1alpha1/tank/statistics/1
```

# Running the Dockerfile

## Build the image

```
docker build -f ./Dockerfile -t trackmyfish .
```

## Run the image

```
docker run -p 8443:8443 -v /path/to/config:/config trackmyfish
```

## Publish the docker image

```
docker login

docker tag trackmyfish simondrake/trackmyfish:v1alpha1

docker push simondrake/trackmyfish:v1alpha1
```

# Tests

Tests are separated into two separate categories

## Unit Tests

Unit Tests do not require any set-up steps and can simply be run with `make test`

## Integration Tests

Integration Tests require actual services to be up (e.g. Postgres), so require a bit of set-up before they can be run.

* Run the test database (`docker-compose -f docker-compose-integration-tests.yaml up -d`)
* Run the `docker-compose` command again, because the DB won't be available when the migrations run (see [#1](https://github.com/TrackMyFish/TrackMyFish/issues/1))
* Run `make integration-test`

**Note:** If you need to login to the db container, you can do so with the following command:

```
docker exec -it integration_tests_db psql postgresql://user:password@localhost:5432/trackmyfishtests
```

# ToDo

* [ ] Write Tests
