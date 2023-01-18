# Golang CQRS blog API

## Architectural Design

<img width="589" alt="Screenshot 2023-01-18 at 21 30 23" src="https://user-images.githubusercontent.com/40946917/213205778-a052af4b-13a6-4e58-a815-d3ecbce15661.png">


The application adopts a clean architectural design with some modifications and CQRS (Command and Query Responsibility Segregation) to address search and query filtering of readable and good data by using elasticsearch. Query results will also be cached for faster request response.

Domains are foldered, for example article domain :
```bash
├── article.go
├── delivery
│   └── http
│       └── handlers
│           └── article.go
├── repositories
│   ├── cache_redis.go
│   ├── elasticsearch.go
│   └── writedb_postgree.go
├── tests
│   ├── http_test.go
│   ├── infrastructure_test.go
│   ├── main_test.go
│   ├── repository_cache_test.go
│   ├── repository_command_test.go
│   ├── repository_query_test.go
│   ├── usecase_command_test.go
│   └── usecase_query_test.go
└── usecases
    ├── cqrs_command.go
    └── cqrs_query.go
```
- this application uses a simple CQRS implementation with separate read and write data stores :
    
    - PostgreeSQl as write and source truth
    - ElasticSearch as read database also as search engine
    - Redis as Cache

    <br>

- every new article data created will be indexed to elastic search (sync). Because it's still a monolith for now it's not really necessary to use event sourcing (async)
- app cache every query with request param for fast access, cache need to be invalidated if new article is created to save source relevant data via (database write)

## How to test
This application uses unit testing to carry out tests, for easier integration testing this application uses dockertest. Once you've run your tests, dockertest spins up all the dependencies/infrastructure needed to run your tests and cleanup when done, so tests are tested on real infrastructure not mocks.

To run integration test use this command :

With make
```bash
make test
```

Without make
```bash
go test ./article/tests/... -v
```

If you just want to run http test endpoint use this command :

```bash
go test ./article/tests/... -run TestHttpFindArticles -v
```

```bash
go test ./article/tests/... -run TestHttpApiCreateArticle  -v
```

## How to run

### Using Docker Compose
1. take a look to docker-compose.yml you can constumize or use default config.
2. when using docker envs loaded from docker-compose.yml or docker env
3. make sure docker running
4. run by using following command

    ```bash
    docker-compose up -d
    ```
5. Run seeder if needed by using command:

    ```bash
    docker exec go-cqrs-api ./app-cli seed
    ```

### Using Docker for Infrastructure only
1. take a look to docker-compose.yml you can constumize or use default config.
2. copy .env.example to .env and configure, or just leave it as default.
3. make sure docker running
4. run infras container by running following command :

    ```bash
    docker-compose up -d go-cqrs-postgree go-cqrs-elasticsearch go-cqrs-redis
    ```
5. migrate database by using command :

    With make
    ```bash
    make migrate
    ```

    Without make
    ```bash
    go run ./cmd/cli/... migrate up
    ```
6. run seeder if needed :

    With make
    ```bash
    make seed
    ```

    Without make
    ```bash
    go run ./cmd/cli/... seed
    ```
7. run api by running command :

    With make :
    ```bash
    make run
    ```
    
    Without make :
    ```bash
    go run ./cmd/api/main.go
    ```


## API Docs

### CURL
Create Article :
```bash
curl -X POST http://localhost:3000/articles \
 -H 'accept: application/json' \
 -H 'Content-Type: application/json' \
 -d '{ 
        "author": "alfiankan", 
        "body": "my blog is very very", 
        "title": "my fisrt" 
    }'
```

Find/Get Article :
```bash
curl -X GET 'http://localhost:3000/articles?keyword=alfiankan&author=cqrs&page=1' \
 -H 'accept: application/json'
```

Find/Get query params:
|param|type|required|note|
|---|---|---|---|
| keyword | string | false | search keyword on title or body |
| author | string | false | filter by author name |
| page | int | false | page result do tou large amount data, page from 1..n every page hold 50 articles, default page is 1 |


### Swagger OAS
- you can access swagger web on `http://<host>:<port>/swagger/index.html`
![screencapture-localhost-3000-swagger-index-html-2023-01-18-22_04_37](https://user-images.githubusercontent.com/40946917/213206702-8cb11691-9af2-4ee9-a532-deb41530b09a.png)


### Postman Documenter
- you can use Postman Documenter by click button below :

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/24530299-d36635f1-b08e-4f61-a85f-8d1f18eecfd9?action=collection%2Ffork&collection-url=entityId%3D24530299-d36635f1-b08e-4f61-a85f-8d1f18eecfd9%26entityType%3Dcollection%26workspaceId%3D554dcc4f-cf17-4e8a-bdb2-bcda713286cf)

## Test and Benchmark

test running 1000 request with 10 concurrent users.

```bash
hey -n 1000 -c 10 'http://localhost:3000/articles?keyword=machine&author=Adam%20Geitgey&page=1'
```

```bash
Summary:
  Total:        0.8775 secs
  Slowest:      0.0679 secs
  Fastest:      0.0048 secs
  Average:      0.0085 secs
  Requests/sec: 1139.6366
  

Response time histogram:
  0.005 [1]     |
  0.011 [930]   |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.017 [56]    |■■
  0.024 [3]     |
  0.030 [0]     |
  0.036 [0]     |
  0.043 [0]     |
  0.049 [0]     |
  0.055 [1]     |
  0.062 [3]     |
  0.068 [6]     |


Latency distribution:
  10% in 0.0061 secs
  25% in 0.0069 secs
  50% in 0.0076 secs
  75% in 0.0087 secs
  90% in 0.0105 secs
  95% in 0.0118 secs
  99% in 0.0547 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0048 secs, 0.0679 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0004 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0016 secs
  resp wait:    0.0065 secs, 0.0035 secs, 0.0643 secs
  resp read:    0.0020 secs, 0.0009 secs, 0.0084 secs

Status code distribution:
  [200] 1000 responses

```


