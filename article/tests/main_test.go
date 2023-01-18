package tests

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/alfiankan/go-cqrs-blog/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestMain(m *testing.M) {

	// set config
	cfg := config.Load("../../.env")

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// POSTGREESQL SETUP
	postgreePorts := []docker.PortBinding{{HostPort: "5432"}}
	pool.RemoveContainerByName("go-cqrs-postgree")

	_, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "go-cqrs-postgree",
		Repository:   "postgres",
		Tag:          "14.1-alpine",
		PortBindings: map[docker.Port][]docker.PortBinding{"5432/tcp": postgreePorts},
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"listen_addresses = '*'",
		}})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// ping postgree before continue
	for {
		pgConn, err := infrastructure.NewPgConnection(cfg)
		if err == nil {
			break
		}
		if err := pgConn.Ping(); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Println("postgree up")

	// ELASTICSEARCH SETUP
	elasticsearchPorts := []docker.PortBinding{{HostPort: "9200"}}
	pool.RemoveContainerByName("go-cqrs-elasticsearch")

	_, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "go-cqrs-elasticsearch",
		Repository:   "docker.elastic.co/elasticsearch/elasticsearch",
		Tag:          "8.6.0",
		PortBindings: map[docker.Port][]docker.PortBinding{"9200/tcp": elasticsearchPorts},
		CapAdd:       []string{"IPC_LOCK"},
		Env: []string{
			"ELASTIC_USERNAME=elastic",
			"ELASTIC_PASSWORD=elastic",
			"xpack.security.enabled=true",
			"discovery.type=single-node",
		}})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// ping wlasticsearch before continue
	for {
		esConn, err := infrastructure.NewElasticSearchClient(cfg)
		if err == nil {
			break
		}
		if _, err := esConn.Cat.Health(); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Println("elasticsearch up")

	// REDIS SETUP
	redisPorts := []docker.PortBinding{{HostPort: "6379"}}
	pool.RemoveContainerByName("go-cqrs-redis")

	_, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "go-cqrs-redis",
		Repository:   "redis",
		Tag:          "6.2-alpine",
		PortBindings: map[docker.Port][]docker.PortBinding{"6379/tcp": redisPorts},
		Cmd:          []string{"--requirepass", "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81"},
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// ping redis before continue
	for {
		redisConn, err := infrastructure.NewRedisConnection(cfg)
		if err == nil {
			break
		}
		if err := redisConn.Ping().Err(); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Println("redis up")

	// SETUP MIGRATION AND SEED
	pgConn, err := infrastructure.NewPgConnection(cfg)
	driver, err := postgres.WithInstance(pgConn, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres",
		driver,
	)
	if err := migrator.Up(); err != nil {
		log.Fatal("migration failed")
	}

	if err := migrations.Seed(); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	pool.RemoveContainerByName("go-cqrs-redis")
	pool.RemoveContainerByName("go-cqrs-elasticsearch")
	pool.RemoveContainerByName("go-cqrs-postgree")

	os.Exit(code)
}
