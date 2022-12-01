package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	Databaser
	types.Copuser
	comfig.Listenerer
	comfig.Logger
}

type config struct {
	getter kv.Getter
	Databaser
	types.Copuser
	comfig.Listenerer
	comfig.Logger
}

func New(getter kv.Getter) Config {
	return config{
		getter:     getter,
		Databaser:  NewDatabaser(getter),
		Copuser:    copus.NewCopuser(getter),
		Listenerer: comfig.NewListenerer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
	}
}
