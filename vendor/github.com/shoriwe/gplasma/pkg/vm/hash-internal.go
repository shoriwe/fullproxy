package vm

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	NotHashable = fmt.Errorf("not hashable")
)

func createHashString(a any) string {
	return fmt.Sprintf("%s -- %v", reflect.TypeOf(a), a)
}

type Hash struct {
	mutex       *sync.Mutex
	internalMap map[string]*Value
}

func (h *Hash) Size() int64 {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return int64(len(h.internalMap))
}

func (h *Hash) Set(key, value *Value) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	switch key.TypeId() {
	case StringId:
		h.internalMap[createHashString(string(key.GetBytes()))] = value
	case BytesId:
		h.internalMap[createHashString(key.GetBytes())] = value
	case BoolId:
		h.internalMap[createHashString(key.GetBool())] = value
	case IntId:
		h.internalMap[createHashString(key.GetInt64())] = value
	case FloatId:
		h.internalMap[createHashString(key.GetFloat64())] = value
	default:
		return NotHashable
	}
	return nil
}

func (h *Hash) Get(key *Value) (*Value, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	switch key.TypeId() {
	case StringId:
		return h.internalMap[createHashString(string(key.GetBytes()))], nil
	case BytesId:
		return h.internalMap[createHashString(key.GetBytes())], nil
	case BoolId:
		return h.internalMap[createHashString(key.GetBool())], nil
	case IntId:
		return h.internalMap[createHashString(key.GetInt64())], nil
	case FloatId:
		return h.internalMap[createHashString(key.GetFloat64())], nil
	}
	return nil, NotHashable
}

func (h *Hash) Del(key *Value) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	switch key.TypeId() {
	case StringId:
		delete(h.internalMap, createHashString(string(key.GetBytes())))
	case BytesId:
		delete(h.internalMap, createHashString(key.GetBytes()))
	case BoolId:
		delete(h.internalMap, createHashString(key.GetBool()))
	case IntId:
		delete(h.internalMap, createHashString(key.GetInt64()))
	case FloatId:
		delete(h.internalMap, createHashString(key.GetFloat64()))
	default:
		return NotHashable
	}
	return nil
}

func (h *Hash) Copy() *Hash {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	result := &Hash{
		mutex:       &sync.Mutex{},
		internalMap: make(map[string]*Value, len(h.internalMap)),
	}
	for key, value := range h.internalMap {
		result.internalMap[key] = value
	}
	return result
}

func (h *Hash) In(key *Value) (bool, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	var found bool
	switch key.TypeId() {
	case StringId:
		_, found = h.internalMap[createHashString(string(key.GetBytes()))]
	case BytesId:
		_, found = h.internalMap[createHashString(key.GetBytes())]
	case BoolId:
		_, found = h.internalMap[createHashString(key.GetBool())]
	case IntId:
		_, found = h.internalMap[createHashString(key.GetInt64())]
	case FloatId:
		_, found = h.internalMap[createHashString(key.GetFloat64())]
	default:
		return false, NotHashable
	}
	return found, nil
}

func (plasma *Plasma) NewInternalHash() *Hash {
	return &Hash{
		mutex:       &sync.Mutex{},
		internalMap: map[string]*Value{},
	}
}
