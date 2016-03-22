package lynx

import (
	"encoding/base64"
	"reflect"
	"testing"
)

var (
	key        = []byte("key")
	nonce      = "He+7w+GWRC3JDWKpgwLYF+9V0WtieRvk"
	cypherText = []byte("4eRv2E8AHA4K6o/k6R5Dx+NMGW2Zlk3Fh8C6c2LwcmWQvsuo97Qi")
	plainText  = []byte("this is some plain text")
)

func TestEncrypt(t *testing.T) {
	encryptedCypherText := Encrypt(plainText, nonce, key)
	if !reflect.DeepEqual(encryptedCypherText, cypherText) {
		t.Errorf("Cypher text %s did not match expected %s", encryptedCypherText, cypherText)
	}
}

func TestDecrypt(t *testing.T) {
	decryptedPlainText, err := Decrypt(cypherText, nonce, key)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(plainText, decryptedPlainText) {
		t.Errorf("Plain text '%s' did not match expected '%s'", decryptedPlainText, plainText)
	}
}

func TestNonce(t *testing.T) {
	nonceBase64, err := NewNonce()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	nonce, _ := base64.StdEncoding.DecodeString(string(nonceBase64))
	if len(nonce) != nonceLen {
		t.Errorf("Unexpected nonce length of %d", len(nonce))
	}
}
