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
	key   string
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	listItem, exists := c.items[key]

	if !exists {
		if c.queue.Len() == c.capacity {
			lastListItem := c.queue.Back()
			cacheItem := lastListItem.Value.(*cacheItem)
			c.queue.Remove(lastListItem)
			delete(c.items, Key(cacheItem.key))
		}
	} else {
		c.queue.Remove(listItem)
	}

	c.items[key] = c.queue.PushFront(&cacheItem{
		key:   string(key),
		value: value,
	})

	return exists
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	var value interface{}
	listItem, exists := c.items[key]

	if exists {
		cacheItem := listItem.Value.(*cacheItem)
		c.queue.MoveToFront(listItem)
		value = cacheItem.value
	}

	return value, exists
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
