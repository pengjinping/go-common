package cache

import (
	"container/list"
	"sync"
	"time"
)

type MemoryStore struct {
	list       map[string]*list.Element //真实存储
	expireList *list.List               //要过期的键(双向链表)
	mu         sync.RWMutex             //读写锁保证并发读写
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		list:       make(map[string]*list.Element),
		expireList: list.New(),
	}
}

func (ms *MemoryStore) Set(key string, value interface{}, time int) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	m := &Memory{}
	m.Set(key, value, time)
	ms.list[key] = ms.expireList.PushBack(m)
}

func (ms *MemoryStore) Forever(key string, value interface{}) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	m := &Memory{}
	m.Forever(key, value)

	listItem := list.New()
	ms.list[key] = listItem.PushBack(m)
}

func (ms *MemoryStore) Get(key string) interface{} {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	value, ok := ms.list[key]
	if !ok {
		return nil
	}

	expire := value.Value.(*Memory).TTL()
	if b := !expire.IsZero() && time.Now().After(expire); b {
		ms.expireList.Remove(ms.list[key])
		delete(ms.list, key)
		return nil
	}

	return value.Value.(*Memory).Get()
}

func (ms *MemoryStore) IsExpire(key string) bool {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// 如果不存在值 则肯定过期
	value, ok := ms.list[key]
	if !ok {
		return true
	}

	expire := value.Value.(*Memory).TTL()
	return !expire.IsZero() && time.Now().After(expire)
}

func (ms *MemoryStore) Has(key string) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if _, ok := ms.list[key]; ok {
		return true
	} else {
		return false
	}
}

func (ms *MemoryStore) Delete(key string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.list[key]; ok {
		ms.expireList.Remove(ms.list[key])
		delete(ms.list, key)
	}
}

type Memory struct {
	key    string
	value  interface{}
	expire time.Time
}

func (m *Memory) Get() interface{} {
	return m.value
}
func (m *Memory) Set(key string, value interface{}, d int) {
	m.key = key
	m.value = value
	m.expire = time.Now().Add(time.Duration(d) * time.Second)
}
func (m *Memory) Forever(key string, value interface{}) {
	m.key = key
	m.value = value
}

func (m *Memory) TTL() time.Time {
	return m.expire
}
