package redisx

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	Init()
	assert.NotNil(t, client)
}

func TestSet(t *testing.T) {
	Init()
	err := Set("test", "test")
	assert.Nil(t, err)
}
func TestHSet(t *testing.T) {
	Init()
	err := Del("test")
	assert.Nil(t, err)
	err = HSet("test", "test", "test", "test2", "test3")
	assert.Nil(t, err)
}

func TestGetString(t *testing.T) {
	Init()
	err := Set("test", "test")
	assert.Nil(t, err)
	val, err := GetString("test")
	assert.Nil(t, err)
	assert.Equal(t, "test", val)
}

func TestGetBool(t *testing.T) {
	Init()
	err := Set("test", "true")
	assert.Nil(t, err)
	val, err := GetBool("test")
	assert.Nil(t, err)
	assert.Equal(t, true, val)
}

func TestGetByte(t *testing.T) {
	Init()
	err := Set("test", "test")
	assert.Nil(t, err)
	val, err := GetByte("test")
	assert.Nil(t, err)
	assert.Equal(t, []byte("test"), val)
}

func TestGetFloat32(t *testing.T) {
	Init()
	err := Set("test", "1.1")
	assert.Nil(t, err)
	val, err := GetFloat32("test")
	assert.Nil(t, err)
	assert.Equal(t, float32(1.1), val)
}
func TestGetFloat64(t *testing.T) {
	Init()
	err := Set("test", "1.1")
	assert.Nil(t, err)
	val, err := GetFloat64("test")
	assert.Nil(t, err)
	assert.Equal(t, 1.1, val)
}

func TestGetInt(t *testing.T) {
	Init()
	err := Set("test", "1")
	assert.Nil(t, err)
	val, err := GetInt("test")
	assert.Nil(t, err)
	assert.Equal(t, 1, val)
}

func TestGetInt64(t *testing.T) {
	Init()
	err := Set("test", "1")
	assert.Nil(t, err)
	val, err := GetInt64("test")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), val)
}

func TestGetUint64(t *testing.T) {
	Init()
	err := Set("test", "1")
	assert.Nil(t, err)
	val, err := GetUint64("test")
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), val)
}

func TestSetNx(t *testing.T) {
	Init()
	err := Del("test")
	assert.Nil(t, err)
	err = SetNx("test", "test", time.Second*20)
	assert.Nil(t, err)
	val, err := GetString("test")
	assert.Nil(t, err)
	assert.Equal(t, "test", val)
}

func TestTTL(t *testing.T) {
	Init()
	err := Del("test")
	assert.Nil(t, err)
	err = SetNx("test", "test", time.Second*20)
	assert.Nil(t, err)
	ttl, err := TTL("test")
	assert.Nil(t, err)
	assert.Equal(t, true, ttl.Seconds() > 19)
}
func TestPersist(t *testing.T) {
	Init()
	err := Del("test")
	assert.Nil(t, err)
	err = SetNx("test", "test", time.Second*20)
	assert.Nil(t, err)
	err = Persist("test")
	assert.Nil(t, err)
	ttl, err := TTL("test")
	assert.Nil(t, err)
	assert.Equal(t, true, ttl.String() == "-1ns")
}
func TestGetStringMapString(t *testing.T) {
	Init()
	err := Del("test")
	assert.Nil(t, err)
	err = HSet("test", "test", "test", "test2", "test3")
	assert.Nil(t, err)
	val, err := GetStringMapString("test")
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"test": "test", "test2": "test3"}, val)
}
func TestGetString2(t *testing.T) {
	Init()
	err := Del("test")
	assert.Nil(t, err)
	str, err := GetString("test2")
	assert.NotNil(t, err)
	assert.Equal(t, "", str)
}

func TestUseConn(t *testing.T) {
	Init()
	err := Del("test")
	assert.Nil(t, err)
	err = UseConn(func(conn *redis.Conn) error {
		return conn.Set(context.Background(), "test2", "test", time.Second*20).Err()
	})
	assert.Nil(t, err)
	val, err := GetString("test2")
	assert.Nil(t, err)
	assert.Equal(t, "test", val)
}
