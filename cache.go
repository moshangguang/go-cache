package gocache

import "context"

type Cache interface {
	Get(ctx context.Context, key string) (data interface{}, ok bool)
	Put(ctx context.Context, key string, data interface{})
	Delete(ctx context.Context, key string) (data interface{})
	Clear(ctx context.Context)
	Len() int
}

type WrapperCache func(Cache) Cache
