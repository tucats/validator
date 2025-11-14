package validator

import (
	"encoding/json"
	"strings"
	"sync"
)

// Dictionary is the map that retains named validators. When a validator
// is defined of type struct, it is stored in the Dictionary. This allows
// recursive references to a validator to be resolved by a dictionary lookup.
// Because validator definitions are thread-safe, a mutex lock serializes
// access to the dictionary.
var Dictionary map[string]*Item = make(map[string]*Item)
var dictionaryLock sync.Mutex

// Store a validator to the Dictionary. This method serializes access to
// the dictionary.
func store(name string, item *Item) error {
	if name == "" || item == nil {
		return nil
	}

	dictionaryLock.Lock()
	defer dictionaryLock.Unlock()

	Dictionary[name] = item

	return nil
}

// Find a named validator from the Dictionary. This method serializes access to
// the dictionary. If the name does not exist in the Dictionary, it returns nil and
// false.
func find(name string) (*Item, bool) {
	dictionaryLock.Lock()
	defer dictionaryLock.Unlock()

	item, found := Dictionary[name]

	return item, found
}

// Define a new validator by name. The validator is created using normal
// reflection, and the resulting validator is stored in the Dictionary.
// This method serializes access to the dictionary. If the name already
// exists in the Dictionary, an error is returned.
func Define(name string, obj any) error {
	if _, found := find(name); found {
		return ErrNameAlreadyExists.Value(name)
	}

	// You cannot define an item with the reserved prefix we used
	// for recursive aliases to existing names.
	if strings.HasPrefix(name, aliasPrefix) {
		return ErrInvalidName.Value(name)
	}

	// Create a new validator from the given object.
	item, err := New(obj)
	if err != nil {
		return err
	}

	// Store the validator in the Dictionary.
	return store(name, item)
}

// DumpJSON is a diagnostic function that allows the user of the
// validator package to dump the contents of the Dictionary in JSON format.
func DumpJSON() []byte {
	b, err := json.MarshalIndent(Dictionary, "", "   ")
	if err != nil {
		return nil
	}

	return b
}
