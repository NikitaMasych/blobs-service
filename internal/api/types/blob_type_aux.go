package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

var ErrBlobTypeInvalid = errors.New("BlobType is invalid")

var _BlobTypeNameToValue = map[string]BlobType{
	"ordinary": BlobTypeOrdinary,
}

var _BlobTypeValueToName = map[BlobType]string{
	BlobTypeOrdinary: "ordinary",
}

func (r BlobType) String() string {
	s, ok := _BlobTypeValueToName[r]
	if !ok {
		return fmt.Sprintf("BlobType(%d)", r)
	}
	return s
}

func (r BlobType) Validate() error {
	_, ok := _BlobTypeValueToName[r]
	if !ok {
		return ErrBlobTypeInvalid
	}
	return nil
}

func (r BlobType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _BlobTypeValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid BlobType: %d", r)
	}
	return json.Marshal(s)
}

func (r *BlobType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("BlobType should be a string, got %s", data)
	}
	v, ok := _BlobTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid BlobType %q", s)
	}
	*r = v
	return nil
}

func (t *BlobType) Scan(src interface{}) error {
	i, ok := src.(int64)
	if !ok {
		return fmt.Errorf("can't scan from %T", src)
	}
	*t = BlobType(i)
	return nil
}

func (t BlobType) Value() (driver.Value, error) {
	return int64(t), nil
}
