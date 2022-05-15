package iface

type Cache interface {
	Get(key string) ([]byte, error)

	Set(key string, value []byte) error

	Delete(key string) error
}
