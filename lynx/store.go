package lynx

import (
	"errors"

	"github.com/snikch/api/ctx"
)

type contextKey int

const (
	storeKey contextKey = iota
)

// Store instances can be used to register Lockable instances against for later
// unlocking or locking.
type Store struct {
	Items   []Lockable
	Context *ctx.Context
}

// ContextStore gets or creates a store for the supplied context.
func ContextStore(context *ctx.Context) *Store {
	if store, ok := context.GetOk(storeKey); ok {
		return store.(*Store)
	}
	store := NewStore(context)
	context.Set(storeKey, store)
	return store
}

// NewStore returns an initialized Store instance.
func NewStore(context *ctx.Context) *Store {
	return &Store{
		Items:   []Lockable{},
		Context: context,
	}
}

// Save adds additional Lockable instances to the store.
func (store *Store) Save(items ...Lockable) {
	store.Items = append(store.Items, items...)
}

// KeyHandler is the function that returns an unlock key for a context.
var KeyHandler func(*ctx.Context) ([]byte, error)

// ErrNoKeyHandler is returned when attempting to find a key with no handler set.
var ErrNoKeyHandler = errors.New("No lynx.KeyHandler function has been set")

// ContextKey returns the key for the supplied context by calling a registered handler.
func ContextKey(context *ctx.Context) ([]byte, error) {
	if KeyHandler == nil {
		return nil, ErrNoKeyHandler
	}
	return KeyHandler(context)
}

// Unlock calls the Unlock method on all Lockable items.
func (store *Store) Unlock() error {
	// Shortcut for empty items slice.
	if len(store.Items) == 0 {
		return nil
	}

	// Get the unlock key for this context.
	key, err := ContextKey(store.Context)
	if err != nil {
		return err
	}
	// Create required variables and start the unlock process.
	var firstErr error
	ch := make(chan error)
	for i := range store.Items {
		go store.unlockIndex(ch, i, key)
	}

	// While it may not be true for some users, we're assuming nobody has added to
	// store.Items between loops.
	for _ = range store.Items {
		err := <-ch
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// unlockIndex will unlock a single item at the supplied index with any error
// being returned on the supplied error channel.
func (store *Store) unlockIndex(ch chan error, i int, key []byte) {
	// Get the fields and nonce
	fields := store.Items[i].LockableValues()

	ch <- UnlockMany(key, fields)
}
