package config

import (
	"net/url"
	"reflect"

	"github.com/spf13/cast"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	horizon "gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
	"gitlab.com/tokend/keypair/figurekeypair"
)

type Horizoner interface {
	Horizon() *horizon.Connector
}

type horizoner struct {
	getter kv.Getter
	once   comfig.Once
	value  *horizon.Connector
	err    error
}

func NewHorizoner(getter kv.Getter) Horizoner {
	return &horizoner{getter: getter}
}

func (h *horizoner) Horizon() *horizon.Connector {
	return h.once.Do(func() interface{} {
		var config struct {
			Endpoint *url.URL     `fig:"endpoint,required"`
			Signer   keypair.Full `fig:"signer,required"`
		}

		err := figure.
			Out(&config).
			With(figure.BaseHooks, figurekeypair.Hooks, URLHook).
			From(kv.MustGetStringMap(h.getter, "horizon")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out horizon"))
		}

		hrz := horizon.NewConnector(config.Endpoint)
		if config.Signer != nil {
			hrz = hrz.WithSigner(config.Signer)
		}
		return hrz
	}).(*horizon.Connector)
}

var URLHook = figure.Hooks{
	"*url.URL": func(value interface{}) (reflect.Value, error) {
		str, err := cast.ToStringE(value)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to parse string")
		}
		u, err := url.Parse(str)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to parse url")
		}
		return reflect.ValueOf(u), nil
	},
}
