/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type CancelDeferredPaymentCreationRequestOp struct {
	Key
	Relationships CancelDeferredPaymentCreationRequestOpRelationships `json:"relationships"`
}
type CancelDeferredPaymentCreationRequestOpResponse struct {
	Data     CancelDeferredPaymentCreationRequestOp `json:"data"`
	Included Included                               `json:"included"`
}

type CancelDeferredPaymentCreationRequestOpListResponse struct {
	Data     []CancelDeferredPaymentCreationRequestOp `json:"data"`
	Included Included                                 `json:"included"`
	Links    *Links                                   `json:"links"`
	Meta     json.RawMessage                          `json:"meta,omitempty"`
}

func (r *CancelDeferredPaymentCreationRequestOpListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *CancelDeferredPaymentCreationRequestOpListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustCancelDeferredPaymentCreationRequestOp - returns CancelDeferredPaymentCreationRequestOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCancelDeferredPaymentCreationRequestOp(key Key) *CancelDeferredPaymentCreationRequestOp {
	var cancelDeferredPaymentCreationRequestOp CancelDeferredPaymentCreationRequestOp
	if c.tryFindEntry(key, &cancelDeferredPaymentCreationRequestOp) {
		return &cancelDeferredPaymentCreationRequestOp
	}
	return nil
}
