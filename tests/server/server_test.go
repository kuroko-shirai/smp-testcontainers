package server

import (
	"context"
	"lab/internal/config"
	"lab/internal/services/srv"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "lab/pkg/proto"

	"github.com/testcontainers/testcontainers-go"
)

type Suite struct {
	suite.Suite

	assert     *assert.Assertions
	controller *gomock.Controller
	redis      testcontainers.Container
	client     pb.PublisherClient
}

func (this *Suite) SetupTest() {
	this.assert = assert.New(this.T())
	this.controller = gomock.NewController(this.T())

	ctx := context.Background()

	redis, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "redis:alpine",
				ExposedPorts: []string{"6379/tcp"},
			},
			Started: true,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	ip, err := redis.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}
	port, err := redis.MappedPort(ctx, "6379")
	if err != nil {
		log.Fatal(err)
	}
	this.redis = redis

	srv, err := srv.New(ctx, &config.Config{
		Server: config.ServerConfig{
			Address: "localhost:8080",
		},
		Redis: config.RedisConfig{
			Address:  ip + ":" + port.Port(),
			Password: "",
			DB:       0,
		},
	})

	go srv.Run(ctx)

	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	this.client = pb.NewPublisherClient(conn)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (this *Suite) TearDownTest() {
	this.controller.Finish()
}

func (this *Suite) TestYourMethod() {
	ctx := context.Background()

	this.redis.Exec(ctx, []string{"redis-cli", "SET", "key", "ok"})

	resp, err := this.client.SendMessage(ctx, &pb.Message{
		Text: "hello",
	})
	if err != nil {
		log.Fatal(err)
	}

	this.assert.Equal("ok", resp.Result)
}
