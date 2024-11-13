package core

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestSyncMaps(t *testing.T) {
	m := NewSyncMap[int, string]()

	m.Set(1, "1")
	v, ok := m.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "1", v)

	v, ok = m.Get(2)
	assert.False(t, ok)
	assert.Empty(t, v)

	exist := m.ContainsKey(1)
	assert.True(t, exist)

	exist = m.ContainsKey(2)
	assert.False(t, exist)

	key, ok := m.Find("1", comparator)

	assert.True(t, ok)
	assert.Equal(t, 1, *key)

	key, ok = m.Find("2", comparator)

	assert.False(t, ok)
	assert.Nil(t, key)

	m.AddRange([]string{"2", "3", "4"}, keygen)

	items := m.Range()
	assert.Equal(t, 4, len(items))
	assert.Contains(t, items, "1", "2", "3", "4")

	m.Remove(1)

	exist = m.ContainsKey(1)
	assert.False(t, exist)

	_, ok = m.TryAdd("2", keygen2, comparator)
	assert.False(t, ok)
}

func keygen(s string) int {
	k, _ := strconv.Atoi(s)
	return k
}

func keygen2() int {
	return 2
}

func comparator(s1 string, s2 string) bool {
	return s1 == s2
}
