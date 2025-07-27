package mod_db

// RedisUnmarshaler is the interface implemented by types that can unmarshal a Redisearch document field.
// The UnmarshalRedis method is called with the raw value from the Redisearch document.
type RedisUnmarshaler interface {
	UnmarshalRedis(value interface{}) error
}
