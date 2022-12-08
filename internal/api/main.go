package api

import (
	"blobs/internal/api/ctx"
	"blobs/internal/api/handlers"
	"blobs/internal/config"
	"blobs/internal/database"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func Serve(cfg config.Config) error {
	r := newRouter(cfg)
	if err := cfg.Copus().RegisterChi(r); err != nil {
		return errors.Wrap(err, "failed to register chi router")
	}

	cfg.Log().Info("Api started on", cfg.Listener().Addr())
	return http.Serve(cfg.Listener(), r)
}

func newRouter(cfg config.Config) chi.Router {
	r := chi.NewRouter()
	r = attachMiddleware(r, cfg)
	return initRoutes(r)
}

func attachMiddleware(m *chi.Mux, cfg config.Config) *chi.Mux {
	m.Use(
		ape.RecoverMiddleware(cfg.Log()),
		ape.LoganMiddleware(cfg.Log()),
		ape.CtxMiddleware(
			ctx.SetLog(cfg.Log()),
			ctx.SetBlobQ(database.NewBlobsQ(cfg.DB())),
		),
	)
	return m
}

func initRoutes(m *chi.Mux) *chi.Mux {
	m.Route("/blobs", func(r chi.Router) {
		m.Post("/", handlers.CreateBlob)
		m.Get("/", handlers.GetBlobs)
		m.Get("/{blob}", handlers.GetBlob)
		m.Delete("/{blob}", handlers.DeleteBlob)
	})
	return m
}
