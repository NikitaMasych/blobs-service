package service

import (
	"context"

	"blobs/internal/api"
	"blobs/internal/config"
	"blobs/internal/service/runners"
)

type Service struct {
	cfg config.Config
}

func NewService(cfg config.Config) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) Run() error {
	cancellable, cancel := context.WithCancel(context.Background())
	defer cancel()

	go runners.NewAssetCreator(s.cfg).Run(cancellable)
	go runners.NewAssetRemover(s.cfg).Run(cancellable)

	return api.Serve(s.cfg)
}
