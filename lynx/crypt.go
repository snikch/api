package lynx

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

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

// MaxNonceAttempts represents how many times we will retry generating a nonce.
var MaxNonceAttempts = 5

// NewNonce promises to return a new nonce for use in encrypting strings.
func NewNonce() (string, error) {
	// Generate a random nonce
	nonce := &[nonceLen]byte{}
	var err error
	// Read the nonce from rand, with retry.
	for i := 0; i < MaxNonceAttempts; i++ {
		if _, err = rand.Read(nonce[:]); err == nil {
			break
		}
	}
	if err != nil {
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

// UnlockMany decrypts an array of cyphers.
func UnlockMany(key []byte, cypherTexts []*string) error {
	ch := make(chan error)

	count := 0
	for _, field := range cypherTexts {
		// Ignore nil pointers.
		if field == nil {
			continue
		}
		count++
		go unlockOne(ch, key, field)
	}

	// Pull off the channel as many times as we pushed onto it.
	var mainErr error
	for i := 0; i < count; i++ {
		err := <-ch
		if mainErr == nil && err != nil {
			mainErr = err
		}
	}
	return mainErr
}

// unlockOne is a goroutine safe way to encrypt a single field.
func unlockOne(ch chan<- error, key []byte, cypherText *string) {
	components := strings.Split(*cypherText, payloadBase64Delimiter)
	if len(components) < 2 {
		ch <- fmt.Errorf("Invalid number of components in cypher text")
		return
	}

	// Decrypt the field.
	plainText, err := Decrypt([]byte(strings.Join(components[1:], payloadBase64Delimiter)), components[0], key)
	if err != nil {
		ch <- err
		return
	}

	// Replace the string at the pointer with the cypherText.
	*cypherText = string(plainText)
	ch <- nil
}

// LockMany decrypts an array of cyphers.
func LockMany(key []byte, plainTexts []*string) error {
	ch := make(chan error)

	count := 0
	for _, field := range plainTexts {
		// Ignore nil pointers.
		if field == nil {
			continue
		}
		count++
		go lockOne(ch, key, field)
	}

	// Pull off the channel as many times as we pushed onto it.
	var mainErr error
	for i := 0; i < count; i++ {
		err := <-ch
		if mainErr == nil && err != nil {
			mainErr = err
		}
	}
	return mainErr
}

// LockString decrypts an array of cyphers.
func LockString(key []byte, plainText *string) error {
	ch := make(chan error)
	go lockOne(ch, key, plainText)
	return <-ch
}

// lockOne is a goroutine safe way to encrypt a single field.
func lockOne(ch chan<- error, key []byte, plainText *string) {
	nonce, err := NewNonce()
	if err != nil {
		ch <- err
		return
	}
	// Encrypt the field.
	cypherText := Encrypt([]byte(*plainText), nonce, key)

	// Replace the string at the pointer with the cypherText.
	*plainText = nonce + payloadBase64Delimiter + string(cypherText)
	ch <- nil
}
