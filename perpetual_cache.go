package gocache

import (
	"context"
)

type PerpetualCache struct {
	data map[string]interface{}
}

func (cache *PerpetualCache) Get(_ context.Context, key string) (data interface{}, ok bool) {
	data, ok = cache.data[key]
	return
}

func (cache *PerpetualCache) Put(_ context.Context, key string, data interface{}) {
	cache.data[key] = data
}

func (cache *PerpetualCache) Delete(ctx context.Context, key string) (data interface{}) {
	data, ok := cache.Get(ctx, key)
	if !ok {
		return
	}
	delete(cache.data, key)
	return
}

func (cache *PerpetualCache) Clear(_ context.Context) {
	cache.data = make(map[string]interface{})
}

func (cache *PerpetualCache) Len() int {
	return len(cache.data)
}

func NewPerpetualCache() Cache {
	return &PerpetualCache{
		data: map[string]interface{}{},
	}
}
