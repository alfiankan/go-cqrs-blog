# Golang CQRS blog API

## Architectural Design

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
    
    1. PostgreeSQl as write and source truth
    2. ElasticSearch as read database also as search engine
    3. Redis as Cache

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
    docker-compose up go-cqrs-postgree go-cqrs-elasticsearch go-cqrs-redis -d
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

### Postman Documenter
- you can use Postman Documenter at

## Test and Benchmark




