package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	ErrBlobTypeInvalid = errors.New("BlobType is invalid")

	_BlobTypeNameToValue = map[string]BlobType{
		"ordinary": BlobTypeOrdinary,
	}

	_BlobTypeValueToName = map[BlobType]string{
		BlobTypeOrdinary: "ordinary",
	}
)

func GetBlobType(v string) (b BlobType, err error) {
	err = b.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, v)))
	if err != nil {
		return b, errors.Wrap(err, "failed to unmarshal blob type")
	}
	return b, nil
}

func (b *BlobType) String() string {
	s, ok := _BlobTypeValueToName[*b]
	if !ok {
		return fmt.Sprintf("BlobType(%d)", b)
	}
	return s
}

func (b *BlobType) Validate() error {
	_, ok := _BlobTypeValueToName[*b]
	if !ok {
		return ErrBlobTypeInvalid
	}
	return nil
}

func (b *BlobType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(b).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _BlobTypeValueToName[*b]
	if !ok {
		return nil, fmt.Errorf("invalid BlobType: %d", b)
	}
	return json.Marshal(s)
}

func (b *BlobType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("BlobType should be a string, got %s", data)
	}
	v, ok := _BlobTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid BlobType %q", s)
	}
	*b = v
	return nil
}

func (b *BlobType) Scan(src interface{}) error {
	i, ok := src.(int64)
	if !ok {
		return fmt.Errorf("can't scan from %T", src)
	}
	*b = BlobType(i)
	return nil
}

func (b *BlobType) Value() (driver.Value, error) {
	return int64(*b), nil
}
