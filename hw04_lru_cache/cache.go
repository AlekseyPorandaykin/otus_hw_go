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
	_, itemExist := cache.items[key]

	if !itemExist && cache.queue.Len() >= cache.capacity {
		lastItem := cache.queue.Back()
		for key, value := range cache.items {
			if value.Value == lastItem.Value {
				delete(cache.items, key)
				break
			}
		}
		cache.queue.Remove(lastItem)
	}
	cache.items[key] = cache.queue.PushFront(value)

	return itemExist
}

func (cache lruCache) Get(key Key) (interface{}, bool) {
	if item, itemExist := cache.items[key]; itemExist {
		cache.queue.MoveToFront(item)
		return item.Value, true
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
