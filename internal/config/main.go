package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/connectors/builder"
	"gitlab.com/tokend/connectors/keyer"
	"gitlab.com/tokend/connectors/submit"
)

type Config interface {
	pgdb.Databaser
	submit.Submission
	types.Copuser
	comfig.Listenerer
	builder.Builderer
	keyer.Keyer
	comfig.Logger
}

type config struct {
	getter kv.Getter
	pgdb.Databaser
	submit.Submission
	types.Copuser
	comfig.Listenerer
	builder.Builderer
	keyer.Keyer
	comfig.Logger
}

func New(getter kv.Getter) Config {
	return config{
		getter:     getter,
		Databaser:  pgdb.NewDatabaser(getter),
		Submission: submit.NewSubmission(getter),
		Copuser:    copus.NewCopuser(getter),
		Listenerer: comfig.NewListenerer(getter),
		Builderer:  builder.NewBuilderer(getter),
		Keyer:      keyer.NewKeyer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
	}
}
