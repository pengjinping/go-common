package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type MemoryStore struct {
	list       map[string]*list.Element //真实存储
	expireList *list.List               //要过期的键(双向链表)
	mu         sync.RWMutex             //读写锁保证并发读写
}

var MemoryName = "Memory" //store的名称复制代码实例化store以及初始化注册cache

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		list:       make(map[string]*list.Element),
		expireList: list.New(),
	}
}

func (ms *MemoryStore) GetStoreName() string {
	return MemoryName
}

func (ms *MemoryStore) Set(key string, value interface{}, time int) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	m := &Memory{}
	m.Set(key, value, time)
	element := ms.expireList.PushBack(m)

	ms.list[key] = element
	return nil
}

func (ms *MemoryStore) Forever(key string, value interface{}) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	m := &Memory{}
	m.Forever(key, value) //无过期时间的只存在list map
	list := list.New()
	element := list.PushBack(m)

	ms.list[key] = element
	return nil
}

func (ms *MemoryStore) Get(key string) (interface{}, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	isExpire, err := ms.IsExpire(key)
	if err != nil {
		return nil, err
	} else {
		if isExpire {
			ms.expireList.Remove(ms.list[key])
			delete(ms.list, key)

			return nil, fmt.Errorf("Cache %s not found", key)
		} else {
			return ms.list[key].Value.(*Memory).Get(), nil
		}
	}
}

func (ms *MemoryStore) Delete(key string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.list[key]; ok {
		ms.expireList.Remove(ms.list[key])
		delete(ms.list, key)
		return nil
	}
	return nil
}

func (ms *MemoryStore) Has(key string) (bool, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if _, ok := ms.list[key]; ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (ms *MemoryStore) Clear() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.list = make(map[string]*list.Element)

	ms.expireList = list.New()
	return nil
}

func (ms *MemoryStore) IsExpire(key string) (bool, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	b, err := ms.Has(key)
	if err != nil {
		return false, err
	}

	if !b {
		return false, err
	}

	expire := ms.list[key].Value.(*Memory).TTL(key)
	return !expire.IsZero() && time.Now().After(expire), nil
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
func (m *Memory) TTL(key string) time.Time {
	return m.expire
}
