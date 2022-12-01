package service

import (
	"net"
	"net/http"

	postgres "blobs/internal/api/database"
	"blobs/internal/config"

	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log      *logan.Entry
	db       *postgres.Repo
	copus    types.Copus
	listener net.Listener
}

func (s *service) run() error {
	s.log.Info("Service started")
	r := s.router(postgres.NewBlobs(s.db))

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func newService(cfg config.Config) *service {
	return &service{
		log:      cfg.Log(),
		db:       cfg.Database(),
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(); err != nil {
		panic(err)
	}
}
