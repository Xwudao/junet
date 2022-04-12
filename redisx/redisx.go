package redisx

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var (
	config = Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
		Ctx:      context.Background(),
	}

	client *redis.Client
)

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int

	Ctx context.Context
}

type Opt func(c *Config)

func SetHost(host string) Opt {
	return func(c *Config) {
		c.Host = host
	}
}

func SetPort(port int) Opt {
	return func(c *Config) {
		c.Port = port
	}
}

func SetPassword(pass string) Opt {
	return func(c *Config) {
		c.Password = pass
	}
}

func SetDb(db int) Opt {
	return func(c *Config) {
		c.DB = db
	}
}

func SetCtx(ctx context.Context) Opt {
	return func(c *Config) {
		c.Ctx = ctx
	}
}

func checkInit() {
	if client == nil {
		panic("redisx not init")
	}
}

func Init(opts ...Opt) {
	for _, opt := range opts {
		opt(&config)
	}

	client = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
}
func UseConn(run func(conn *redis.Conn) error) error {
	conn := client.Conn(config.Ctx)
	err := run(conn)
	_ = conn.Close()
	return err
}

func GetClient() *redis.Client {
	checkInit()
	return client
}

func Set(key string, val any) error {
	checkInit()
	return client.Set(config.Ctx, key, val, 0).Err()
}
func SetNx(key string, val any, expiration time.Duration) error {
	checkInit()
	return client.SetNX(config.Ctx, key, val, expiration).Err()
}

// HSet accepts values in following formats:
//  - HSet("myhash", "key1", "value1", "key2", "value2")
//  - HSet("myhash", []string{"key1", "value1", "key2", "value2"})
//  - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
func HSet(key string, val ...any) error {
	checkInit()
	return client.HSet(config.Ctx, key, val...).Err()
}

func Del(key string) error {
	checkInit()
	return client.Del(config.Ctx, key).Err()
}

func TTL(key string) (time.Duration, error) {
	checkInit()
	return client.TTL(config.Ctx, key).Result()
}

func Persist(key string) error {
	checkInit()
	return client.Persist(config.Ctx, key).Err()
}

func GetString(key string) (string, error) {
	checkInit()
	return client.Get(config.Ctx, key).Result()
}

func GetInt(key string) (int, error) {
	checkInit()
	return client.Get(config.Ctx, key).Int()
}
func GetBool(key string) (bool, error) {
	checkInit()
	return client.Get(config.Ctx, key).Bool()
}
func GetFloat64(key string) (float64, error) {
	checkInit()
	return client.Get(config.Ctx, key).Float64()
}
func GetInt64(key string) (int64, error) {
	checkInit()
	return client.Get(config.Ctx, key).Int64()
}
func GetUint64(key string) (uint64, error) {
	checkInit()
	return client.Get(config.Ctx, key).Uint64()
}
func GetFloat32(key string) (float32, error) {
	checkInit()
	return client.Get(config.Ctx, key).Float32()
}
func GetByte(key string) ([]byte, error) {
	checkInit()
	return client.Get(config.Ctx, key).Bytes()
}

func GetStringMapString(key string) (map[string]string, error) {
	checkInit()
	return client.HGetAll(config.Ctx, key).Result()
}
