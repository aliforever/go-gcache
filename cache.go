package cache

import (
	"sync"
	"time"
)

type Entry[T any] struct {
	Item     *T
	ExpireAt time.Time
}

type Cache[T any] struct {
	lock sync.Mutex

	items map[string]Entry[T]
}

func New[T any]() *Cache[T] {
	return &Cache[T]{
		items: make(map[string]Entry[T]),
	}
}

func (c *Cache[T]) Load(key string) (*T, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if val, ok := c.items[key]; ok && !val.ExpireAt.IsZero() && val.ExpireAt.After(time.Now()) {
		return val.Item, true
	}

	return nil, false
}

func (c *Cache[T]) LoadVal(key string) (T, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if val, ok := c.items[key]; ok && !val.ExpireAt.IsZero() && val.ExpireAt.After(time.Now()) {
		return *val.Item, true
	}

	var val T

	return val, false
}

func (c *Cache[T]) LoadOrStoreFunc(key string, fetchFunc func() (*T, error), ttl time.Duration) (*T, bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if val, ok := c.items[key]; ok && !val.ExpireAt.IsZero() && val.ExpireAt.After(time.Now()) {
		return val.Item, true, nil
	}

	item, err := fetchFunc()
	if err != nil {
		return nil, false, err
	}

	c.items[key] = Entry[T]{
		Item:     item,
		ExpireAt: time.Now().Add(ttl),
	}

	return item, false, nil
}

func (c *Cache[T]) LoadOrStoreValFunc(key string, fetchFunc func() (T, error), ttl time.Duration) (T, bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if val, ok := c.items[key]; ok && !val.ExpireAt.IsZero() && val.ExpireAt.After(time.Now()) {
		return *val.Item, true, nil
	}

	item, err := fetchFunc()
	if err != nil {
		var v T

		return v, false, err
	}

	c.items[key] = Entry[T]{
		Item:     &item,
		ExpireAt: time.Now().Add(ttl),
	}

	return item, false, nil
}

func (c *Cache[T]) Store(key string, data *T, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items[key] = Entry[T]{
		Item:     data,
		ExpireAt: time.Now().Add(ttl),
	}
}

func (c *Cache[T]) StoreVal(key string, data T, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items[key] = Entry[T]{
		Item:     &data,
		ExpireAt: time.Now().Add(ttl),
	}
}
