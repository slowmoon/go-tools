package consisthash

import (
	"sort"
	"strconv"
)

type Hash  func([]byte) uint32

type Map struct {
	hash   Hash
    replicate int
    keys []int
    keyMaps map[int]string
}

func New(f Hash, replicate int) *Map {
	return  &Map{
		hash:      f,
		replicate: replicate,
		keyMaps:   make(map[int]string),
	}
}

func (m *Map)IsEmpty() bool {
    return len(m.keys) ==0
}

func (m *Map)Add(keys...string)  {
	for _, key := range keys {
		for i:=0;i < m.replicate; i ++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			if _, ok := m.keyMaps[hash];!ok {
				m.keys = append(m.keys, hash)
				m.keyMaps[hash]  = key
			}
		}
	}
	sort.Ints(m.keys)
}

func (m *Map)Get(key string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.hash([]byte(key)))

	s := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] > hash
	})

	if s == len(m.keys)	{
		s = 0
	}
	return m.keyMaps[s]
}





