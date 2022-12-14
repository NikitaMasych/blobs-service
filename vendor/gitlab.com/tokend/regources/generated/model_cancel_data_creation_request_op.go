/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type CancelDataCreationRequestOp struct {
	Key
	Relationships CancelDataCreationRequestOpRelationships `json:"relationships"`
}
type CancelDataCreationRequestOpResponse struct {
	Data     CancelDataCreationRequestOp `json:"data"`
	Included Included                    `json:"included"`
}

type CancelDataCreationRequestOpListResponse struct {
	Data     []CancelDataCreationRequestOp `json:"data"`
	Included Included                      `json:"included"`
	Links    *Links                        `json:"links"`
	Meta     json.RawMessage               `json:"meta,omitempty"`
}

func (r *CancelDataCreationRequestOpListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *CancelDataCreationRequestOpListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustCancelDataCreationRequestOp - returns CancelDataCreationRequestOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCancelDataCreationRequestOp(key Key) *CancelDataCreationRequestOp {
	var cancelDataCreationRequestOp CancelDataCreationRequestOp
	if c.tryFindEntry(key, &cancelDataCreationRequestOp) {
		return &cancelDataCreationRequestOp
	}
	return nil
}
