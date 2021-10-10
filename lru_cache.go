package gocache

import (
	"context"
	"time"
)

type LRUCache struct {
	lruList  *LruList
	cap      int
	delegate Cache
}

func (cache *LRUCache) Get(ctx context.Context, key string) (data interface{}, ok bool) {
	data, ok = cache.delegate.Get(ctx, key)
	if !ok {
		return
	}
	cache.lruList.Get(key)
	return
}

func (cache *LRUCache) Put(ctx context.Context, key string, data interface{}) {
	if len(key) == 0 {
		panic("Not allow key length is zero")
	}
	if data == nil {
		panic("Not allow data is nil")
	}
	cache.delegate.Put(ctx, key, data)
}

func (cache *LRUCache) Delete(ctx context.Context, key string) (data interface{}) {
	data = cache.delegate.Delete(ctx, key)
	if data == nil {
		return
	}
	cache.lruList.Delete(key)
	return
}

func (cache *LRUCache) Clear(ctx context.Context) {
	cache.delegate.Clear(ctx)
}

func (cache *LRUCache) Len() int {
	return cache.delegate.Len()
}

func NewLRUCache(delegate Cache) Cache {
	return &LRUCache{
		lruList: &LruList{
			m:             map[string]*LruNode{},
			head:          nil,
			tail:          nil,
			evictDuration: time.Minute * 10,
		},
		delegate: delegate,
	}
}

type LruList struct {
	m             map[string]*LruNode
	head, tail    *LruNode
	evictDuration time.Duration
}

func (lru *LruList) Delete(key string) *LruNode {
	node, ok := lru.m[key]
	if !ok {
		return nil
	}
	if lru.OnlyOne(node) {
		lru.Reset()
		return node.Reset()
	}
	switch node {
	case lru.head:
		lru.head = node.next
		lru.head.prev = nil
	case lru.tail:
	default:

	}
	return node
}
func (lru *LruList) OnlyOne(node *LruNode) bool {
	return lru.head == lru.tail && lru.head == node
}
func (lru *LruList) Put(key string) (evict []string) {
	return
}
func (lru *LruList) Get(key string) (node *LruNode, ok bool) {
	node, ok = lru.m[key]
	if !ok {
		return
	}
	node.DelayAlive(lru.evictDuration)
	if lru.OnlyOne(node) {
		return
	}
	switch node {
	case lru.head:
		lru.head = node.next
		lru.head.prev = nil
		lru.tail.next = node
		node.prev = lru.tail
		lru.tail = node
	case lru.tail:
	default:
		prev := node.prev
		next := node.next
		prev.next = next
		next.prev = prev
		lru.tail.next = node
		node.prev = lru.tail
		lru.tail = node
	}
	return
}
func (lru *LruList) Reset() {
	lru.head, lru.tail = nil, nil
	lru.m = map[string]*LruNode{}
}

type LruNode struct {
	prev, next *LruNode
	key        string
	alive      int64
}

func (node *LruNode) Reset() *LruNode {
	node.prev, node.next = nil, nil
	return node
}

func (node *LruNode) DelayAlive(evict time.Duration) {
	node.alive = time.Now().Add(evict).UnixNano()
}
