package lynx

import "sync"

// Lockable defines an interface for objects that can be locked and unlocked.
type Lockable interface {
	LockableValues() []*string
}

// Lock takes a lockable and encrypts its fields with the supplied key. If no
// nonce is set, one will be set for it.
func Lock(key []byte, lockable Lockable) error {
	fields := lockable.LockableValues()

	wg := &sync.WaitGroup{}
	for _, field := range fields {
		// We accept nil values, but do not attempt to encrypt them.
		if field == nil {
			continue
		}
		wg.Add(1)
		go lockField(wg, key, field)
	}

	wg.Wait()
	return nil
}

// lockField is a goroutine safe way to encrypt a single field.
func lockField(wg *sync.WaitGroup, key []byte, plainText *string) {
	nonce, err := NewNonce()
	if err != nil {
		panic(err)
	}
	// Encrypt the field
	cypherText := Encrypt([]byte(*plainText), nonce, key)

	// Replace the string at the pointer with the cypherText.
	*plainText = nonce + ":" + string(cypherText)
	wg.Done()
}

// Unlock takes a lockable and encrypts its fields with the supplied key.
func Unlock(key []byte, lockable Lockable) error {
	fields := lockable.LockableValues()
	return UnlockMany(key, fields)
}
