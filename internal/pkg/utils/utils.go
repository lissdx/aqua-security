package utils

import (
	"errors"
	"fmt"
)

// Map maps the in[]
// if f func(S) T returns transformed value
// returns transformed array or nil if in[] is nil
func Map[S, T any](f func(S) T, in []S) []T {
	if in == nil {
		return nil
	}
	result := make([]T, 0, len(in))
	for i := 0; i < len(in); i++ {
		v := f(in[i])
		result = append(result, v)
	}
	return result
}

// ObjectMap maps the single object
// if f func(*S) *T returns transformed value
// returns transformed object or nil if in is nil
func ObjectMap[S, T any](f func(*S) *T, in *S) *T {
	if in == nil {
		return nil
	}

	return f(in)
}

// Filter filtered the in[] array by boolean
// if f func(T) bool returns FALSE the value will be removed
// returns filtered array or nil if in[] is nil
func Filter[T any](f func(T) bool, in []T) []T {
	if in == nil {
		return nil
	}
	result := make([]T, 0, len(in))

	for i := 0; i < len(in); i++ {
		if f(in[i]) {
			v := in[i]
			result = append(result, v)
		}
	}

	return result
}

func IsNulPointer[T any](p *T) bool {
	return p == nil
}

// Coalesce returns firs found Not nil pointer
// if all args are nil pointers the nil
// will be returned
func Coalesce[T any](args ...*T) *T {
	for i := 0; i < len(args); i++ {
		if !IsNulPointer(args[i]) {
			return args[i]
		}
	}
	return nil
}

// MapAny maps any object to the given type
// returns pointer to object or error if given object is nil or there is no way
// to convert the object
func MapAny[E any](v any) (*E, error) {
	switch t := v.(type) {
	case nil:
		return nil, errors.New("cannot map given nil object")
	case E:
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot map given object type: %T to type: %T", v, t)
	}
}

func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}
	return result
}

// RemoveDuplicates got a comparable array
// and removes duplicated values
// if nil or an empty array is provided
// the empty array will be returned
// the original (provided) array will not be changed
// the result will be stored in new array
// the original order is not guaranteed
func RemoveDuplicates[E comparable](l []E) []E {
	if len(l) <= 0 {
		res := make([]E, 0)
		return res
	}

	hMap := make(map[E]bool)

	for i := 0; i < len(l); i++ {
		if !hMap[l[i]] {
			hMap[l[i]] = true
		}
	}

	res := make([]E, 0, len(hMap))
	for k, _ := range hMap {
		res = append(res, k)
	}
	return res
}
