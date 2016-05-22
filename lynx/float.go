package lynx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// EncryptedFloat represents a float that is encryptable as a string.
type EncryptedFloat struct {
	str *string
}

// NewEncryptedFloat returns a newly initialized EncryptedFloat.
func NewEncryptedFloat(f float64) EncryptedFloat {
	ef := EncryptedFloat{}
	ef.SetFloat(f)
	return ef
}

// SetFloat updates the underlying float value.
func (ef *EncryptedFloat) SetFloat(f float64) {
	s := fmt.Sprintf("%g", f)
	ef.str = &s
}

// Float64 returns a float64 representation of the EncryptedFloat.
func (ef *EncryptedFloat) Float64() (float64, error) {
	if ef.str == nil {
		return 0, fmt.Errorf("EncryptedFloat.Float64: attempting to return nil float64")
	}
	return strconv.ParseFloat(*ef.str, 64)
}

// UnmarshalJSON will unmarshall a raw json representation of a number into an EncryptedFloat.
func (ef *EncryptedFloat) UnmarshalJSON(raw []byte) error {
	f := 0.0
	// Ensure that the value is a valid json number.
	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}
	// Use the unmarshaled float value.
	s := fmt.Sprintf("%g", f)
	ef.str = &s
	return nil
}

// MarshalJSON marshals the current value into a json number.
func (ef EncryptedFloat) MarshalJSON() ([]byte, error) {
	if ef.str == nil {
		return []byte("0.0"), nil
	}
	// Convert to float if possible.
	f, err := strconv.ParseFloat(*ef.str, 64)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert encrypted float to float: %s", err)
	}
	// Then json marshal the float value.
	return json.Marshal(f)
}

// EncryptableString returns the underlying string pointer.
func (ef EncryptedFloat) EncryptableString() *string {
	return ef.str
}

// String implements the stringer interface.
func (ef EncryptedFloat) String() string {
	if ef.str == nil {
		return ""
	}
	return *ef.str
}

// Scan implements sql.Scanner for scanning database values.
func (ef *EncryptedFloat) Scan(value interface{}) error {
	switch val := value.(type) {
	case float64:
		s := fmt.Sprintf("%f", val)
		ef.str = &s
	case int64:
		s := fmt.Sprintf("%d", val)
		ef.str = &s
	case []byte:
		s := string(val)
		ef.str = &s
	case string:
		ef.str = &val
	default:
		return fmt.Errorf("EncryptedFloat.Scan: invalid scan type %T", value)
	}
	return nil
}

// Value implements value.Valuer to provide database values.
func (ef EncryptedFloat) Value() (driver.Value, error) {
	return ef.str, nil
}
