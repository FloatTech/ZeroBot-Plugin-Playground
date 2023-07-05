package store

import "sync"

type onlineKv = map[string]string

var (
	onlineMu    sync.RWMutex
	onlineStore = make(map[string][]onlineKv)
)

func CacheOnline(uid string, messages []onlineKv) {
	onlineMu.Lock()
	defer onlineMu.Unlock()
	onlineStore[uid] = messages
}

func DeleteOnline(uid string) {
	onlineMu.Lock()
	defer onlineMu.Unlock()
	delete(onlineStore, uid)
}

func GetOnline(uid string) []onlineKv {
	if result, ok := onlineStore[uid]; ok {
		return result
	}
	return make([]onlineKv, 0)
}
