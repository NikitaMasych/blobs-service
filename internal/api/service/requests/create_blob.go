package requests

import (
	"encoding/base32"
	"encoding/json"
	"fmt"
	"net/http"

	"blobs/internal/api/resources"
	"blobs/internal/api/types"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/hash"
)

func NewCreateBlobRequest(r *http.Request) (resources.BlobRequest, error) {
	var request resources.BlobRequestResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request.Data, errors.Wrap(err, "failed to unmarshal")
	}

	return request.Data, ValidateCreateBlobRequest(request.Data)
}

func ValidateCreateBlobRequest(r resources.BlobRequest) error {
	return validation.Errors{
		"/data/type":             validation.Validate(&r.Type, validation.Required),
		"/data/attributes/value": validation.Validate(&r.Attributes.Value, validation.Required),
	}.Filter()
}

func Blob(r resources.BlobRequest) (*types.Blob, error) {
	blob, err := types.GetBlobType(string(r.Type))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create blob")
	}
	msg := fmt.Sprintf("%d%s", blob, r.Attributes.Value)
	hash := hash.Hash([]byte(msg))

	return &types.Blob{
		ID:    base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(hash[:]),
		Type:  blob,
		Value: r.Attributes.Value,
	}, nil
}
