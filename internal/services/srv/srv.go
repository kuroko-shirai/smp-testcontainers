package srv

import (
	"context"
	"log"
	"net"

	"lab/internal/config"
	pb "lab/pkg/proto"

	redis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedPublisherServer

	storage  *redis.Client
	listener *net.Listener
}

func New(ctx context.Context, cfg *config.Config) (*Server, error) {
	storage := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := storage.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	listener, err := net.Listen(
		"tcp",
		cfg.Server.Address,
	)
	if err != nil {
		return nil, err
	}

	return &Server{
		storage:  storage,
		listener: &listener,
	}, nil
}

func (this *Server) Run(ctx context.Context) error {
	go func() {
		server := grpc.NewServer()

		pb.RegisterPublisherServer(server, this)

		if err := server.Serve(*this.listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	for {
	}
}

func (this *Server) Stop(context.Context) error {
	return this.storage.Close()
}

type Message struct {
	ID      string
	Message string
}

func (this *Server) SendMessage(
	ctx context.Context,
	msg *pb.Message,
) (*pb.Response, error) {
	res := this.storage.Get(context.Background(), "key")

	return &pb.Response{
		Result: res.Val(),
	}, nil
}
