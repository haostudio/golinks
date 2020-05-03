package memory

type store interface {
	get(keys ...string) (value []byte, exists bool)
	set(value []byte, keys ...string)
	del(key ...string)
	iter(f func(key string, val []byte) (next bool), key ...string) (
		exists bool, err error)
	drop(key ...string)
}
