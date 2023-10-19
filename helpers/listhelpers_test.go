package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceContains(t *testing.T) {
	assert.True(t, SliceContains([]string{"asd", "zxc", "qwe"}, "asd"))
	assert.False(t, SliceContains([]string{"asd", "qwe", "zxc"}, "kjg"))
	assert.True(t, SliceContains([]string{"a", "b", "c"}, "a"))
	assert.False(t, SliceContains([]string{"asd", "bnm", "cvb"}, "a"))
}

func TestIsPlainMap(t *testing.T) {
	assert.True(t, IsPlainMap(map[string]interface{}{"asd": 1, "qwe": "a", "zxc": 1.2}))
	assert.False(t, IsPlainMap(map[string]interface{}{"asd": 1, "qwe": map[string]interface{}{"xyz": 2}, "zxc": 1.2}))
}

func TestIsPlainSlice(t *testing.T) {
	assert.True(t, IsPlainSlice([]interface{}{"asd", 1, 2.3}))
	assert.False(t, IsPlainSlice([]interface{}{"asd", 1, map[string]interface{}{"asd": 12, "aqwe": "asddfds"}}))
}

func TestIsNonStringFloatBool(t *testing.T) {
	assert.True(t, IsNonStringFloatBool([]interface{}{"asd", 2}))
	assert.True(t, IsNonStringFloatBool([]string{"asd", "2"}))
	assert.True(t, IsNonStringFloatBool(2))
	assert.False(t, IsNonStringFloatBool(2.3))
	assert.False(t, IsNonStringFloatBool("asd"))
	assert.False(t, IsNonStringFloatBool(true))
}

func TestDiscard(t *testing.T) {
	var newSlice = Discard([]interface{}{1, 2, 3}, 1)
	assert.NotEmpty(t, newSlice)
	assert.True(t, newSlice[0] == 2)
	assert.True(t, newSlice[1] == 3)
	newSlice = Discard([]interface{}{1, 2, 3}, 2)
	assert.NotEmpty(t, newSlice)
	assert.True(t, newSlice[0] == 1)
	assert.True(t, newSlice[1] == 3)
	newSlice = Discard([]interface{}{1, 2, 3}, 3)
	assert.NotEmpty(t, newSlice)
	assert.True(t, newSlice[0] == 1)
	assert.True(t, newSlice[1] == 2)
}

func TestMapContainsNil(t *testing.T) {
	assert.True(t, mapContainsNil(map[string]interface{}{"asd": 1, "qwe": nil}))
	assert.False(t, mapContainsNil(map[string]interface{}{"asd": 1, "qwe": 34.3, "ewrt": "qwer"}))
}

func TestSliceContainsNil(t *testing.T) {
	assert.True(t, sliceContainsNil([]interface{}{1, "asdfs", nil, 2.3}))
	assert.False(t, sliceContainsNil([]interface{}{1, "asdfs", 2.3}))
}
