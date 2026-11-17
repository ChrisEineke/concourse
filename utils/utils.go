package utils

import (
	"reflect"
	"strings"

	"github.com/chigopher/pathlib"
	"github.com/go-viper/mapstructure/v2"
	"github.com/multiformats/go-multiaddr"
)

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Map[T any, U any](input []T, tr func(T) (U, error)) ([]U, error) {
	result := make([]U, len(input))
	for i, v := range input {
		tmp, err := tr(v)
		if err != nil {
			return nil, err
		}
		result[i] = tmp
	}
	return result, nil
}

func StringsToMultiaddrsHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		from reflect.Type,
		to reflect.Type,
		data any,
	) (any, error) {
		if from != reflect.TypeFor[string]() {
			return data, nil
		}
		if to != reflect.TypeFor[[]multiaddr.Multiaddr]() {
			return data, nil
		}
		vals := strings.Split(data.(string), ",")
		return Map(vals, func(str string) (multiaddr.Multiaddr, error) {
			return multiaddr.NewMultiaddr(str)
		})
	}
}

func StringToPathHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		from reflect.Type,
		to reflect.Type,
		data any,
	) (any, error) {
		if from != reflect.TypeFor[string]() {
			return data, nil
		}
		if to != reflect.TypeFor[pathlib.Path]() {
			return data, nil
		}
		return pathlib.NewPath(data.(string)), nil
	}
}
