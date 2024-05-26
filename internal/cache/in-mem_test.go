package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCache_GetDel(t *testing.T) {
	// given
	cache := New()
	cache.Put("key1", "value1", time.Minute)

	// when
	value1, ok1 := cache.GetDel("key1")
	value2, ok2 := cache.GetDel("key1")

	// then
	require.True(t, ok1)
	require.Equal(t, "value1", value1)
	require.False(t, ok2)
	require.Nil(t, value2)
}

func TestCache_Get(t *testing.T) {
	// given
	cache := New()
	cache.Put("key1", "value1", time.Minute)

	// when
	value1, ok1 := cache.Get("key1")
	value2, ok2 := cache.Get("key1")

	// then
	require.True(t, ok1)
	require.Equal(t, "value1", value1)
	require.True(t, ok2)
	require.Equal(t, "value1", value2)
}

func TestNew(t *testing.T) {
	// when
	cache := New()

	// then
	require.NotNil(t, cache)
}

func TestCache_GetDel_Expired(t *testing.T) {
	// given
	cache := New()

	// when
	cache.Put("key1", "value1", time.Nanosecond)
	time.Sleep(10 * time.Nanosecond)
	value, ok := cache.GetDel("key1")

	// then
	require.False(t, ok)
	require.Nil(t, value)
}
