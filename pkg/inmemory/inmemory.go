package inmemory

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

const (
	notFoundID         = 0
	maxStorageCapacity = 1000
)

var (
	errStorageOverloaded = errors.New("storage is overloaded")
	errInvalidObjectID   = errors.New("stirage ibject ID is invalid")
	errObjectNotFound    = errors.New("object not found in storage")
)

type EmptyObject struct{}

func (e *EmptyObject) GetID() int {
	return notFoundID
}

type StorableObject interface {
	GetID() int
	// GetField(fieldName string) interface{}
}

type InMemoryStorage struct {
	mu          *sync.RWMutex
	lastID      atomic.Int64
	size        atomic.Int64
	readStorage []StorableObject
	keyStorage  map[int]int
}

func NewInMemoryStorage(initSize int) *InMemoryStorage {
	return &InMemoryStorage{
		mu:          &sync.RWMutex{},
		lastID:      atomic.Int64{},
		size:        atomic.Int64{},
		readStorage: make([]StorableObject, 0, initSize),
		keyStorage:  make(map[int]int, initSize),
	}
}

func (s *InMemoryStorage) AddObject(ctx context.Context, obj StorableObject) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var newIndex int
	for newIndex = 0; newIndex < len(s.readStorage); newIndex++ {
		if s.readStorage[newIndex].GetID() == notFoundID {
			break
		}
	}

	if newIndex == len(s.readStorage) || len(s.readStorage) == 0 {
		if size := s.size.Load(); size >= maxStorageCapacity {
			return errStorageOverloaded
		}
		s.readStorage = append(s.readStorage, obj)
	} else {
		s.readStorage[newIndex] = obj
	}
	s.keyStorage[obj.GetID()] = newIndex
	s.lastID.Add(1)
	s.size.Add(1)

	return nil
}

func (s *InMemoryStorage) GetObjectByID(ctx context.Context, id int) StorableObject {
	s.mu.RLock()
	defer s.mu.RUnlock()

	index, ok := s.keyStorage[id]
	if !ok {
		return nil
	}

	obj := s.readStorage[index]
	return obj
}

func (s *InMemoryStorage) GetManyObjects(ctx context.Context, offset int, limit int) ([]StorableObject, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	objects := make([]StorableObject, 0, limit)
	var id, count int
	for id = offset + 1; id <= int(s.lastID.Load()) && count < limit; id++ {
		if index, ok := s.keyStorage[id]; ok {
			obj := s.readStorage[index]
			objects = append(objects, obj)
			count++
		}
	}

	return objects, nil
}

func (s *InMemoryStorage) GetAll(ctx context.Context) ([]StorableObject, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	objects := make([]StorableObject, 0, s.size.Load())
	for _, obj := range s.readStorage {
		if obj.GetID() != -1 {
			objects = append(objects, obj)
		}
	}

	return objects, nil
}

func (s *InMemoryStorage) DeleteObject(ctx context.Context, id int) error {
	if id <= notFoundID || id > int(s.lastID.Load()) {
		return errInvalidObjectID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	index, ok := s.keyStorage[id]
	if !ok {
		return errObjectNotFound
	}

	s.readStorage[index] = &EmptyObject{}
	return nil
}
