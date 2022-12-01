/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type RelationCollection struct {
	Data []Key `json:"data"`
}

func (r RelationCollection) MarshalJSON() ([]byte, error) {
	if r.Data == nil {
		r.Data = []Key{}
	}

	type temp RelationCollection
	return json.Marshal(temp(r))
}
