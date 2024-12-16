package infra

import (
	"context"

	"lab/internal/config"
	"lab/internal/services/srv"
)

type Service interface {
	Run(context.Context) error
	Stop(context.Context) error
}

type App struct {
	Services []Service
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	ms, err := srv.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Services: []Service{
			ms,
		},
	}, nil
}

func (a *App) Run(ctx context.Context) {
	for _, s := range a.Services {
		s.Run(ctx)
	}
}

func (a *App) Stop(cxt context.Context) {
	for _, s := range a.Services {
		s.Stop(cxt)
	}
}
