package inmemory

import (
	"errors"
	"sync"
)

type Cache struct {
	mu   *sync.RWMutex
	byId map[string]string
}

var (
	ErrorNotFound = errors.New("Value not found")
)

func New() *Cache {
	return &Cache{
		mu:   &sync.RWMutex{},
		byId: make(map[string]string),
	}
}

func (c *Cache) GetValue(id string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.byId[id]
	if !ok {
		return "", ErrorNotFound
	}
	return value, nil
}

func (c *Cache) CreateOrUpdate(key, value string) (bool, error) {
	var exist bool
	_, err := c.GetValue(key)
	if err == ErrorNotFound {
		exist = true
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.byId[key] = value
	return exist, nil
}
