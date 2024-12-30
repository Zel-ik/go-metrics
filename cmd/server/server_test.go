package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddAndGetList(t *testing.T) {
	memCashe := NewMemoryStore()
	memCashe.add("3", "4")

	assert.NotEmpty(t, memCashe.getList())
	assert.Equal(t, memCashe.metrics["3"], "4")
}

func TestChangeAndGetList(t *testing.T) {
	memCashe := NewMemoryStore()
	memCashe.add("3", "4")
	memCashe.change("3", "6.123")

	assert.NotEmpty(t, memCashe.getList())
	assert.Equal(t, len(memCashe.getList()), 1)
	assert.Equal(t, memCashe.metrics["3"], "6.123")
}
