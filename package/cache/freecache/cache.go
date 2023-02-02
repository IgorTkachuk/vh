package freecache

import (
	"github.com/coocood/freecache"
	"sync"
	"vh/package/cache"
)

type repository struct {
	sync.Mutex
	cache *freecache.Cache
}

func NewCacheRepo(size int) cache.Repository {
	return &repository{
		cache: freecache.NewCache(size),
	}
}

func (r repository) GetIterator() cache.Iterator {
	return &iterator{r.cache.NewIterator()}
}

func (r repository) Get(uuid []byte) ([]byte, error) {
	r.Lock()
	defer r.Unlock()

	return r.cache.Get(uuid)
}

func (r repository) Set(uuid []byte, value []byte, expireIn int) error {
	r.Lock()
	defer r.Unlock()

	return r.cache.Set(uuid, value, expireIn)
}

func (r repository) Del(uuid []byte) (affected bool) {
	r.Lock()
	defer r.Unlock()

	return r.cache.Del(uuid)
}

func (r repository) EntryCount() int64 {
	r.Lock()
	defer r.Unlock()

	return r.cache.EntryCount()
}

func (r repository) HitCount() int64 {
	r.Lock()
	defer r.Unlock()

	return r.cache.HitCount()
}

func (r repository) MissCount() int64 {
	r.Lock()
	defer r.Unlock()

	return r.cache.MissCount()
}
