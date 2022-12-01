package service

import (
	"blobs/internal/api/contracts"
	"blobs/internal/api/service/handlers"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router(blobQ contracts.Blobs) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxBlobQ(blobQ),
		),
	)
	r.Route("/blobs", func(r chi.Router) {
		r.Post("/", handlers.CreateBlob)
		r.Get("/", handlers.GetBlobs)
		r.Get("/{blob}", handlers.GetBlob)
		r.Delete("/{blob}", handlers.DeleteBlob)
	})

	return r
}
