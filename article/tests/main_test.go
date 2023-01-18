package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alfiankan/go-cqrs-blog/common"
	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestMain(m *testing.M) {

	// set override config
	envs := map[string]string{
		"PG_DATABASE_HOST":        "127.0.0.1",
		"PG_DATABASE_USERNAME":    "postgres",
		"PG_DATABASE_PASSWORD":    "postgres",
		"PG_DATABASE_NAME":        "postgres",
		"PG_DATABASE_PORT":        "2345",
		"PG_DATABASE_SSL_MODE":    "disable",
		"LOG_LEVEL":               "debug",
		"ELASTICSEARCH_ADDRESSES": "http://127.0.0.1:2900",
		"ELASTICSEARCH_USERNAME":  "elastic",
		"ELASTICSEARCH_PASSWORD":  "elastic",
		"REDIS_HOST":              "127.0.0.1:9376",
		"REDIS_PASSWORD":          "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		"HTTP_API_PORT":           "5000",
	}

	for key, val := range envs {
		os.Setenv(key, val)
	}
	cfg := config.Load("void")

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// POSTGREESQL SETUP
	postgreePorts := []docker.PortBinding{{HostPort: strconv.Itoa(cfg.PostgreePort)}}
	pool.RemoveContainerByName("go-cqrs-postgree-test")

	if _, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "go-cqrs-postgree-test",
		Repository:   "postgres",
		Tag:          "14.1-alpine",
		PortBindings: map[docker.Port][]docker.PortBinding{"5432/tcp": postgreePorts},
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", cfg.PostgreeUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", cfg.PostgreePass),
			"listen_addresses = '*'",
		}}); err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// ELASTICSEARCH SETUP
	parsedEnv := strings.Split(cfg.ElasticSearchAdresses[0], ":")
	elasticsearchPorts := []docker.PortBinding{{HostPort: parsedEnv[len(parsedEnv)-1]}}
	pool.RemoveContainerByName("go-cqrs-elasticsearch-test")

	if _, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "go-cqrs-elasticsearch-test",
		Repository:   "docker.elastic.co/elasticsearch/elasticsearch",
		Tag:          "8.6.0",
		PortBindings: map[docker.Port][]docker.PortBinding{"9200/tcp": elasticsearchPorts},
		CapAdd:       []string{"IPC_LOCK"},
		Env: []string{
			fmt.Sprintf("ELASTIC_USERNAME=%s", cfg.ElasticSearchUsername),
			fmt.Sprintf("ELASTIC_PASSWORD=%s", cfg.ElasticSearchPassword),
			"xpack.security.enabled=true",
			"discovery.type=single-node",
		}}); err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// REDIS SETUP
	parsedEnv = strings.Split(cfg.RedisHost[0], ":")
	redisPorts := []docker.PortBinding{{HostPort: parsedEnv[len(parsedEnv)-1]}}
	pool.RemoveContainerByName("go-cqrs-redis-test")

	if _, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "go-cqrs-redis-test",
		Repository:   "redis",
		Tag:          "6.2-alpine",
		PortBindings: map[docker.Port][]docker.PortBinding{"6379/tcp": redisPorts},
		Cmd:          []string{"--requirepass", cfg.RedisPass},
	}); err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// PING ALL CHECK CONNECTION OK
	for {
		log.Println("try to ping redis ⏳")
		redisConn, _ := infrastructure.NewRedisConnection(cfg)
		if err := redisConn.Ping(context.Background()).Err(); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	log.Println("redis up and running ✅")

	for {
		log.Println("try to ping postgree ⏳")
		pgConn, _ := infrastructure.NewPgConnection(cfg)
		if err := pgConn.Ping(); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	log.Println("postgree up and running ✅")

	for {
		log.Println("try to ping elasticsearch ⏳")

		esConn, _ := infrastructure.NewElasticSearchClient(cfg)
		if _, err := esConn.Ping(); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	log.Println("elasticsearch up and running ✅")

	// SETUP MIGRATION AND SEED
	if err := common.Migration("../.."); err != nil {
		log.Fatal(err)
	}
	time.Sleep(10 * time.Second)

	if err := common.Seed("../.."); err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	pool.RemoveContainerByName("go-cqrs-redis-test")
	pool.RemoveContainerByName("go-cqrs-elasticsearch-test")
	pool.RemoveContainerByName("go-cqrs-postgree-test")
	os.Exit(code)
}
