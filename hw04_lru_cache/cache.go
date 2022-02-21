package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	item := cacheItem{key: key, value: value}
	_, itemExist := cache.items[item.key]
	if !itemExist && cache.queue.Len() >= cache.capacity {
		lastItem := cache.queue.Back()
		delete(cache.items, lastItem.Value.(cacheItem).key)
		cache.queue.Remove(lastItem)
	}
	cache.items[item.key] = cache.queue.PushFront(item)
	return itemExist
}

func (cache lruCache) Get(key Key) (interface{}, bool) {
	if item, itemExist := cache.items[key]; itemExist {
		cache.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}
	return nil, false
}

func (cache *lruCache) Clear() {
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
