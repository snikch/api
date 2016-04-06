package lynx

import (
	"testing"

	"github.com/snikch/api/ctx"
)

type testUnlocker struct {
	field1 string
	field2 string
	called bool
	nonce  string
}

func init() {
	KeyHandler = func(*ctx.Context) ([]byte, error) {
		return key, nil
	}
}

func (t *testUnlocker) Nonce() string {
	return t.nonce
}

func (t *testUnlocker) SetNonce(nonce string) {
	t.nonce = nonce
}

func (t *testUnlocker) LockableValues() []*string {
	t.called = true
	return []*string{&t.field1, &t.field2}
}

func TestUnlock(t *testing.T) {
	l1 := &testUnlocker{
		field1: nonce + payloadBase64Delimiter + string(cypherText),
		field2: nonce + payloadBase64Delimiter + string(cypherText),
	}
	l2 := &testUnlocker{
		field1: nonce + payloadBase64Delimiter + string(cypherText),
		field2: nonce + payloadBase64Delimiter + string(cypherText),
	}

	context := ctx.NewContext()
	store := NewStore(context)
	store.Save(l1, l2)

	err := store.Unlock()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !l1.called || !l2.called {
		t.Errorf("Unlocker not called")
		return
	}
	expected := string(plainText)
	if !(l1.field1 == expected && l1.field2 == expected && l2.field1 == expected && l2.field2 == expected) {
		t.Errorf("Unexpected field values: %s", l1.field1)
	}
}

func TestSingleError(t *testing.T) {
	newNonce, _ := NewNonce()
	// Supply one correct and one incorrect unlocker.
	l1 := &testUnlocker{
		field1: newNonce + payloadBase64Delimiter + string(cypherText),
		field2: newNonce + payloadBase64Delimiter + string(cypherText),
	}
	l2 := &testUnlocker{
		field1: nonce + payloadBase64Delimiter + string(cypherText),
		field2: nonce + payloadBase64Delimiter + string(cypherText),
	}

	context := ctx.NewContext()
	store := NewStore(context)
	store.Save(l1, l2)

	err := store.Unlock()
	if err == nil {
		t.Errorf("Unexpected lack of an error")
	}
}

func TestMultiError(t *testing.T) {
	newNonce, _ := NewNonce()
	// Both of these unlockers will fail since they have the wrong nonce.
	l1 := &testUnlocker{
		field1: newNonce + payloadBase64Delimiter + string(cypherText),
		field2: newNonce + payloadBase64Delimiter + string(cypherText),
	}
	l2 := &testUnlocker{
		field1: newNonce + payloadBase64Delimiter + string(cypherText),
		field2: newNonce + payloadBase64Delimiter + string(cypherText),
	}

	context := ctx.NewContext()
	store := NewStore(context)
	store.Save(l1, l2)

	err := store.Unlock()
	if err == nil {
		t.Errorf("Unexpected lack of an error")
	}
}
