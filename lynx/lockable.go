package lynx

import "sync"

// Lockable defines an interface for objects that can be locked and unlocked.
type Lockable interface {
	LockableValues() []*string
	Nonce() string
	SetNonce(string)
}

// EnsureNonce will set a nonce on a lockable object if one doesn't exist.
func EnsureNonce(lockable Lockable) error {
	if nonce := lockable.Nonce(); len(nonce) == 0 {
		nonce, err := NewNonce()
		if err != nil {
			return err
		}
		lockable.SetNonce(nonce)
	}
	return nil
}

// Lock takes a lockable and encrypts its fields with the supplied key. If no
// nonce is set, one will be set for it.
func Lock(key []byte, lockable Lockable) error {
	err := EnsureNonce(lockable)
	if err != nil {
		return err
	}

	fields := lockable.LockableValues()
	nonce := lockable.Nonce()

	wg := &sync.WaitGroup{}
	for _, field := range fields {
		// We accept nil values, but do not attempt to encrypt them.
		if field == nil {
			continue
		}
		wg.Add(1)
		go lockField(wg, key, field, nonce)
	}

	wg.Wait()
	return nil
}

// lockField is a goroutine safe way to encrypt a single field.
func lockField(wg *sync.WaitGroup, key []byte, plainText *string, nonce string) {
	// Decrypt the field
	cypherText := Encrypt([]byte(*plainText), nonce, key)

	// Replace the string at the pointer with the cypherText.
	*plainText = string(cypherText)
	wg.Done()
}

// Unlock takes a lockable and encrypts its fields with the supplied key.
func Unlock(key []byte, lockable Lockable) error {
	fields := lockable.LockableValues()
	nonce := lockable.Nonce()
	ch := make(chan error)

	count := 0
	for _, field := range fields {
		// Ignore nil pointers.
		if field == nil {
			continue
		}
		count++
		go unlockField(ch, key, field, nonce)
	}

	// Pull off the channel as many times as we pushed onto it.
	var mainErr error
	for i := 0; i < count; i++ {
		err := <-ch
		if mainErr == nil {
			mainErr = err
		}
	}
	return nil
}

// unlockField is a goroutine safe way to encrypt a single field.
func unlockField(ch chan<- error, key []byte, cypherText *string, nonce string) {
	// Decrypt the field
	plainText, err := Decrypt([]byte(*cypherText), nonce, key)
	if err != nil {
		ch <- err
		return
	}

	// Replace the string at the pointer with the cypherText.
	*cypherText = string(plainText)
	ch <- nil
}
