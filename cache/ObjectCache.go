package cache

type ObjectCache interface {
	Get(key string, item interface{})

	Set(key string, item interface{}) error

	GetString(key string) (string, bool)

	SetString(key string, value string)
}
