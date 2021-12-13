package cache

import (
	"container/list"
	"sync"
	"time"
)

// 要过期的键(双向链表)
var expireList = make(map[string]*list.List)

type MemoryStore struct {
	UUID       string                   //租户名称
	list       map[string]*list.Element //真实存储
	mu         sync.RWMutex             //读写锁保证并发读写
}

func NewMemoryStore(uuid string) *MemoryStore {
	expireList[uuid] = list.New()

	return &MemoryStore{
		UUID:       uuid,
		list:       make(map[string]*list.Element),
	}
}

func (ms *MemoryStore) GetTenant() string {
	return ms.UUID
}

func (ms *MemoryStore) SetTenant(tenant string, tenantId int) bool {
	ms.UUID = tenant

	if _, ok := expireList[tenant]; !ok {
		expireList[tenant] = list.New()
	}

	return true
}

func (ms *MemoryStore) Set(key string, value interface{}, time int) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	key = ms.setKey(key)

	m := &Memory{}
	if time > 0 {
		// 清除以前的值 以最新值存储
		if e := ms.list[key]; e != nil {
			expireList[ms.UUID].Remove(e)
		}

		m.Set(key, value, time)
		ms.list[key] = expireList[ms.UUID].PushBack(m)
	} else {
		m.Forever(key, value)

		listItem := list.New()
		ms.list[key] = listItem.PushBack(m)
	}
}
func (ms *MemoryStore) Forever(key string, value interface{}) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	key = ms.setKey(key)

	m := &Memory{}
	m.Forever(key, value)

	listItem := list.New()
	ms.list[key] = listItem.PushBack(m)
}
func (ms *MemoryStore) Get(key string) interface{} {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	key = ms.setKey(key)

	value, ok := ms.list[key]
	if !ok {
		return nil
	}

	expire := value.Value.(*Memory).TTL()
	if b := !expire.IsZero() && time.Now().After(expire); b {
		expireList[ms.UUID].Remove(ms.list[key])
		delete(ms.list, key)
		return nil
	}

	return value.Value.(*Memory).Get()
}

func (ms *MemoryStore) IsExpire(key string) bool {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	key = ms.setKey(key)

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
	key = ms.setKey(key)

	if _, ok := ms.list[key]; ok {
		return true
	} else {
		return false
	}
}
func (ms *MemoryStore) Delete(key string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	key = ms.setKey(key)

	if _, ok := ms.list[key]; ok {
		expireList[ms.UUID].Remove(ms.list[key])
		delete(ms.list, key)
	}
}

func (ms *MemoryStore) GC() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:			// 触发定时器
			for _, listItem := range expireList {
				l := listItem.Len()
				element := listItem.Front()
				for i := 0; i < l; i++ {
					if element == nil {
						break
					}

					m := element.Value.(*Memory)
					// 先获取到下个元素 防止Remove后无法获取下个元素
					next := element.Next()
					if b := !m.expire.IsZero() && time.Now().After(m.expire); b {
						listItem.Remove(element)
						delete(ms.list, m.key)
					}
					element = next
				}
			}
		}
	}
}

func (ms *MemoryStore) Keys() interface{} {
	var result []string
	listItem := expireList[ms.UUID]
	result = make([]string, listItem.Len())

	i := 0
	for e := listItem.Front(); e != nil; e = e.Next() {
		result[i] = e.Value.(*Memory).key
		i = i + 1
	}

	return result
}

func (ms *MemoryStore) setKey(key string) string {
	return ms.UUID + ":" + key
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
