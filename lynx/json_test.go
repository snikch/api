package lynx

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestJSONUnmarshalEJ(t *testing.T) {
	ej := EncryptedJSON{}
	err := json.Unmarshal([]byte(`{"foo":"bar"}`), &ej)
	if err != nil {
		t.Errorf("unexpected error unmarshalling value: %s", err)
		return
	}
	if ej.String() != `{"foo":"bar"}` {
		t.Errorf("unexpected string value: %s", ej.String())
		return
	}
}

func TestJSONMarshalEJ(t *testing.T) {
	ej := NewEncryptedJSON(`{"foo":"bar"}`)
	m, err := json.Marshal(ej)
	if err != nil {
		t.Errorf("unexpected error marshalling value: %s", err)
		return
	}

	if !reflect.DeepEqual(m, []byte(`{"foo":"bar"}`)) {
		t.Errorf("unexpected marshalled value: %s", string(m))
		return
	}
}

func TestUpdatingStringPointerEJ(t *testing.T) {
	ej := NewEncryptedJSON(`{"foo":"bar"}`)
	s := ej.EncryptableString()
	*s = `{"foo":"bar"}`
	if ej.String() != `{"foo":"bar"}` {
		t.Errorf("unexpected string value: %s", ej.String())
		return
	}
}

func TestScanEJ(t *testing.T) {
	for _, e := range []struct {
		Type   string
		Value  interface{}
		Valid  bool
		Output string
	}{
		{"string", `{"foo":"bar"}`, true, `{"foo":"bar"}`},
		{"[]byte", []byte(`{"foo":"bar"}`), true, `{"foo":"bar"}`},
	} {
		ej := EncryptedJSON{}
		err := ej.Scan(e.Value)
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
	}
}
