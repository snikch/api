package lynx

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/service/kms"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	kmsKeySpec             = kms.DataKeySpecAes256
	nonceLen               = 24
	keyLen                 = 32
	payloadBase64Delimiter = ":"
)

// Key is the type used for keys. This is the main structure used for encrypting/decrypting files and credntials.
type Key struct {
	KeyBase64    string             `json:"key_base64"`
	Key          []byte             `json:"-"`
	DecryptedKey []byte             `json:"-"`
	Context      map[string]*string `json:"context"`
	MasterID     string             `json:"master_id"`
	Region       string             `json:"region"`
}

// Payload is the basic payload container for enc/dec operations.
type Payload struct {
	CypherText []byte
	Nonce      []byte
}

func NewNonce() (string, error) {
	// Generate a random nonce
	nonce := &[nonceLen]byte{}
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(nonce[:])), nil
}

// Encrypt takes a plain text value and returns the encrypted cyphertext.
func Encrypt(plainText []byte, nonceBase64 string, key []byte) []byte {
	plainBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(plainText)))
	base64.StdEncoding.Encode(plainBase64, plainText)

	nonce := make([]byte, base64.StdEncoding.DecodedLen(len(nonceBase64)))
	base64.StdEncoding.Decode(nonce, []byte(nonceBase64))

	nonceArray := &[nonceLen]byte{}
	copy(nonceArray[:], nonce)

	encKey := &[keyLen]byte{}
	copy(encKey[:], key)

	cypherText := []byte{}
	cypherText = secretbox.Seal(cypherText, plainText, nonceArray, encKey)
	cypherBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(cypherText)))
	base64.StdEncoding.Encode(cypherBase64, cypherText)
	return cypherBase64
}

// Decrypt decrypts a payload and returns the decrypted plaintext.
func Decrypt(cypherBase64 []byte, nonceBase64 string, key []byte) ([]byte, error) {
	var plaintext []byte

	nonceText := make([]byte, base64.StdEncoding.DecodedLen(len(nonceBase64)))
	_, err := base64.StdEncoding.Decode(nonceText, []byte(nonceBase64))
	if err != nil {
		return nil, err
	}

	cypherText := make([]byte, base64.StdEncoding.DecodedLen(len(cypherBase64)))
	cypherLen, err := base64.StdEncoding.Decode(cypherText, cypherBase64)
	if err != nil {
		return nil, err
	}

	nonceArray := &[nonceLen]byte{}
	copy(nonceArray[:], nonceText[:nonceLen])

	encKey := &[keyLen]byte{}
	copy(encKey[:], key)

	plaintext, ok := secretbox.Open(plaintext, cypherText[:cypherLen], nonceArray, encKey)
	if !ok {
		return nil, fmt.Errorf("Error decrypting")
	}
	return plaintext, nil
}
