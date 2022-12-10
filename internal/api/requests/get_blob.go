package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
)

type GetBlobRequest struct {
	BlobID string `json:"id"`
}

func NewGetBlobRequest(r *http.Request) (GetBlobRequest, error) {
	request := GetBlobRequest{
		BlobID: chi.URLParam(r, "blob"),
	}
	return request, request.Validate()
}

func (r GetBlobRequest) Validate() error {
	err := Errors{
		"blob": Validate(&r.BlobID, Required),
	}
	return err.Filter()
}
