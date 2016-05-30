package lynx

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestJSONUnmarshal(t *testing.T) {
	ef := EncryptedFloat{}
	err := json.Unmarshal([]byte(`12.50`), &ef)
	if err != nil {
		t.Errorf("unexpected error unmarshalling value: %s", err)
		return
	}
	if ef.String() != "12.5" {
		t.Errorf("unexpected string value: %s", ef.String())
		return
	}
	f, err := ef.Float64()
	if err != nil {
		t.Errorf("unexpected error converting to float64: %s", err)
		return
	}
	if f != 12.50 {
		t.Errorf("unexpected float value: %f", f)
		return
	}
}

func TestJSONMarshal(t *testing.T) {
	ef := NewEncryptedFloat(12.5)
	m, err := json.Marshal(ef)
	if err != nil {
		t.Errorf("unexpected error marshalling value: %s", err)
		return
	}

	if !reflect.DeepEqual(m, []byte(`12.5`)) {
		t.Errorf("unexpected marshalled value: %s", string(m))
		return
	}
}

func TestUpdatingStringPointer(t *testing.T) {
	ef := NewEncryptedFloat(12.50)
	s := ef.EncryptableString()
	*s = "12.3400"
	if ef.String() != "12.3400" {
		t.Errorf("unexpected string value: %s", ef.String())
		return
	}
	f, err := ef.Float64()
	if err != nil {
		t.Errorf("unexpected error getting float64 value: %s", err)
		return
	}

	if f != 12.34 {
		t.Errorf("unexpected float64 value: %f", f)
	}
}

func TestScan(t *testing.T) {
	for _, e := range []struct {
		Type   string
		Value  interface{}
		Valid  bool
		Output float64
	}{
		{"string", "12.50", true, 12.5},
		{"[]byte", []byte("12.5"), true, 12.5},
		{"float64", 12.50, true, 12.5},
		{"nil", nil, true, 0},
		{"int64", int64(300), true, 300},
	} {
		ef := EncryptedFloat{}
		err := ef.Scan(e.Value)
		if e.Valid && err != nil {
			t.Errorf("Failed to scan %s: %s", e.Type, err)
			continue
		} else if !e.Valid && err == nil {
			t.Errorf("Failed to not scan %s", e.Type)
			continue
		}

		if !e.Valid {
			continue
		}

		f, err := ef.Float64()
		if err != nil && e.Type != "nil" {
			t.Errorf("Failed to convert %s to float: %s", e.Type, err)
			continue
		}
		if f != e.Output {
			t.Errorf("Unexpected output for type %s, expected %f got %f", e.Type, e.Output, f)
		}
	}
}

func TestValue(t *testing.T) {
	ef := EncryptedFloat{}
	v, err := ef.Value()
	if err != nil {
		t.Errorf("Unexpected error getting driver.Value: %s", err)
		return
	}

	if !reflect.DeepEqual(nil, v) {
		t.Errorf("Unexpected driver.Value: %s", v)
		return
	}

	ef.SetFloat(12.50)
	v, err = ef.Value()
	if err != nil {
		t.Errorf("Unexpected error getting driver.Value: %s", err)
		return
	}

	s := "12.5"
	if !reflect.DeepEqual(s, v) {
		t.Errorf("Unexpected driver.Value: %s", v)
		return
	}

}
