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
	if time > 0 {
		// 清除以前的值 以最新值存储
		if e := ms.list[key]; e != nil {
			ms.expireList.Remove(e)
		}

		m.Set(key, value, time)
		ms.list[key] = ms.expireList.PushBack(m)
	} else {
		m.Forever(key, value)

		listItem := list.New()
		ms.list[key] = listItem.PushBack(m)
	}
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

func (ms *MemoryStore) GC() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:			// 触发定时器
			l := ms.expireList.Len()
			element := ms.expireList.Front()
			for i := 0; i < l; i++ {
				if element == nil {
					break
				}

				m := element.Value.(*Memory)
				// 先获取到下个元素 防止Remove后无法获取下个元素
				next := element.Next()
				if b := !m.expire.IsZero() && time.Now().After(m.expire); b {
					ms.expireList.Remove(element)
					delete(ms.list, m.key)
				}
				element = next
			}
		}
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
