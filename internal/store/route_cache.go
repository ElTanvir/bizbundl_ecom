package store

type CacheEntry struct {
	Data []byte
	ETag string
}

func SetRouteBytes(path string, data []byte,etag string) {
	entry := CacheEntry{
		Data: data,
		ETag: etag,
	}
	s := Get()
	cacheKey := "route:" + path
	s.Set(cacheKey, entry)
}

func GetRouteBytes(path string) (CacheEntry, bool) {
	s := Get()
	cacheKey := "route:" + path
	v, keyExists := s.Get(cacheKey)
	if !keyExists {
		return CacheEntry{}, false
	}
	u, typeOK := v.(CacheEntry)
	return u, typeOK
}
