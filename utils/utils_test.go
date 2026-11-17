package utils

import (
	"errors"
	"reflect"
	"testing"

	"github.com/chigopher/pathlib"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		result, err := Map(input, func(i int) (string, error) {
			return string(rune(i)), nil
		})
		assert.NoError(t, err, "unexpected error")
		assert.Len(t, result, 0, "expected empty slice, got %v", result)
	})

	t.Run("non-empty slice", func(t *testing.T) {
		input := []int{1, 2, 3}
		expected := []string{"A", "B", "C"}
		result, err := Map(input, func(i int) (string, error) {
			return string(rune(i + 64)), nil
		})
		assert.NoError(t, err, "unexpected error")
		assert.Len(t, result, 3, "expected length 3, got %d", len(result))
		assert.Equal(t, expected, result, "unexpected result: %v", result)
	})

	t.Run("error in transformation", func(t *testing.T) {
		input := []int{1, 2, 3}
		_, err := Map(input, func(i int) (string, error) {
			if i == 2 {
				return "", errors.New("test error")
			}
			return string(rune(i + 64)), nil
		})
		assert.Error(t, err, "expected error")
	})
}

func TestStringsToMultiaddrsHookFunc(t *testing.T) {
	hook := StringsToMultiaddrsHookFunc()

	t.Run("non-string input", func(t *testing.T) {
		input := 123
		result, err := hook(reflect.TypeOf(input), reflect.TypeOf([]multiaddr.Multiaddr{}), input)
		assert.NoError(t, err, "unexpected error: %v", err)
		assert.Equal(t, result, input, "expected input unchanged, got %v", result)
	})

	t.Run("non-multiaddr slice output", func(t *testing.T) {
		input := "test"
		result, err := hook(reflect.TypeOf(input), reflect.TypeOf(""), input)
		assert.NoError(t, err, "unexpected error: %v", err)
		assert.Equal(t, result, input, "expected input unchanged, got %v", result)
	})

	t.Run("valid string input", func(t *testing.T) {
		input := "/ip4/127.0.0.1/tcp/8080,/ip4/127.0.0.1/tcp/9090"
		result, err := hook(reflect.TypeOf(input), reflect.TypeOf([]multiaddr.Multiaddr{}), input)
		assert.NoError(t, err, "unexpected error: %v", err)
		multiaddrs, ok := result.([]multiaddr.Multiaddr)
		assert.True(t, ok, "expected result to be []multiaddr.Multiaddr")
		assert.Len(t, multiaddrs, 2, "expected length 2, got %d", len(multiaddrs))
	})

	t.Run("invalid multiaddr string", func(t *testing.T) {
		input := "/invalid"
		_, err := hook(reflect.TypeOf(input), reflect.TypeOf([]multiaddr.Multiaddr{}), input)
		assert.Error(t, err, "expected error, got nil")
	})
}

func TestStringToPath(t *testing.T) {
	hook := StringToPath()

	t.Run("non-string input", func(t *testing.T) {
		input := 123
		result, err := hook(reflect.TypeOf(input), reflect.TypeOf(pathlib.Path{}), input)
		assert.NoError(t, err, "unexpected error: %v", err)
		assert.Equal(t, input, result, "expected input unchanged, got %v", result)
	})

	t.Run("non-path output", func(t *testing.T) {
		input := "/some/path"
		result, err := hook(reflect.TypeOf(input), reflect.TypeOf(""), input)
		assert.NoError(t, err, "unexpected error: %v", err)
		assert.Equal(t, input, result, "expected input unchanged, got %v", result)
	})

	t.Run("valid string input", func(t *testing.T) {
		input := "/some/path"
		result, err := hook(reflect.TypeOf(input), reflect.TypeOf(pathlib.Path{}), input)
		assert.NoError(t, err, "unexpected error: %v", err)
		path, ok := result.(*pathlib.Path)
		assert.True(t, ok, "expected result to be pathlib.Path")
		assert.Equal(t, input, path.String(), "expected path %s, got %s", input, path.String())
	})
}
