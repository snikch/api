package lynx

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

// EncryptedJSON represents json that is encryptable as a string.
type EncryptedJSON struct {
	str *string
}

// NewEncryptedJSON returns a newly initialized EncryptedJSON.
func NewEncryptedJSON(j string) EncryptedJSON {
	return EncryptedJSON{
		str: &j,
	}
}

// MarshalJSON returns a []byte representation of the string.
func (j EncryptedJSON) MarshalJSON() ([]byte, error) {
	if j.str == nil {
		return []byte("null"), nil
	}
	return []byte(*j.str), nil
}

// UnmarshalJSON stores the supplied []byte data as a string.
func (j *EncryptedJSON) UnmarshalJSON(data []byte) error {
	if data == nil {
		return errors.New("EncryptedJSON: UnmarshalJSON on nil data")
	}
	str := string(data)
	j.str = &str
	return nil
}

// Value implements value.Valuer to provide database values.
func (j EncryptedJSON) Value() (driver.Value, error) {
	if j.str == nil {
		return nil, nil
	}
	return *j.str, nil
}

// Scan implements sql.Scanner for scanning database values.
func (j *EncryptedJSON) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch val := src.(type) {
	case string:
		j.str = &val
	case []byte:
		str := string(val)
		j.str = &str
	default:
		return fmt.Errorf("Incompatible type for EncryptedJSON: %T", src)
	}
	return nil
}

// EncryptableString returns the underlying string pointer.
func (j EncryptedJSON) EncryptableString() *string {
	return j.str
}

// String implements the stringer interface.
func (j EncryptedJSON) String() string {
	if j.str == nil {
		return ""
	}
	return *j.str
}
