package config

import (
	postgres "blobs/internal/api/database"
	"sync"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Databaser interface {
	Database() *postgres.Repo
}

type databaser struct {
	kv kv.Getter
}

func NewDatabaser(kv kv.Getter) Databaser {
	return &databaser{kv}
}

func (d *databaser) Database() *postgres.Repo {
	mu := new(sync.RWMutex)
	mu.Lock()
	defer mu.Unlock()

	config := struct {
		URL     string `fig:"url,required"`
		MaxIdle int    `fig:"max_idle"`
		MaxOpen int    `fig:"max_open"`
	}{
		MaxIdle: 4,
		MaxOpen: 12,
	}

	if err := figure.Out(&config).From(kv.MustGetStringMap(d.kv, "db")).Please(); err != nil {
		panic(errors.Wrap(err, "failed to db"))
	}

	repo, err := postgres.Open(config.URL)
	if err != nil {
		panic(suppressStack(errors.Wrap(err, "failed to open database")))
	}

	repo.DB.SetMaxIdleConns(config.MaxIdle)
	repo.DB.SetMaxOpenConns(config.MaxOpen)

	if err := repo.DB.Ping(); err != nil {
		panic(suppressStack(errors.Wrap(err, "database is not accessible")))
	}

	return repo
}
