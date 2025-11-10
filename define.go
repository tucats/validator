package validator

import (
	"encoding/json"
	"strings"
	"sync"
)

var Dictionary map[string]*Item = make(map[string]*Item)
var dictionaryLock sync.Mutex

func store(name string, item *Item) {
	dictionaryLock.Lock()
	defer dictionaryLock.Unlock()

	Dictionary[name] = item
}

func find(name string) (*Item, bool) {
	dictionaryLock.Lock()
	defer dictionaryLock.Unlock()

	item, found := Dictionary[name]

	return item, found
}

func Define(name string, obj any) error {
	if strings.HasPrefix(name, aliasPrefix) {
		return ErrInvalidName.Value(name)
	}

	item, err := New(obj)
	if err != nil {
		return err
	}

	store(name, item)

	return nil
}

func DumpJSON() []byte {
	b, err := json.MarshalIndent(Dictionary, "", "   ")
	if err != nil {
		return nil
	}

	return b
}
