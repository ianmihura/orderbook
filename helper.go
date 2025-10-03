package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

type f32 = float32
type f64 = float64
type u64 = uint64
type i32 = int32

type Tuple struct {
	a, b, c any
}

type BaseError struct {
	message string
	data    any
}

func (e *BaseError) Error() string {
	return fmt.Sprintf(e.message, e.data)
}

func Assert(t *testing.T, ok bool, message ...any) {
	if !ok {
		t.Error(message...)
	}
}

// func Ok(ok bool, message ...any) {
// 	if !ok {
// 		panic(message)
// 	}
// }

func IsSortedFuncDesc[S ~[]E, E any](x S, cmp func(a, b E) int) bool {
	for i := len(x) - 1; i > 0; i-- {
		if cmp(x[i], x[i-1]) > 0 {
			return false
		}
	}
	return true
}

// Returns x cut to n decimal places.
// *Does not round.*
// Usefull to truncate random prices.
func Truncate(x f64, n int) f64 {
	return math.Floor(x*math.Pow(10, f64(n))) * math.Pow(10, -f64(n))
}

// Custom normal dist variable `X ~ N(mean, std)“.
// Truncated to t decimal places.
func NormFloat32T(mean, std f32, t int) f32 {
	return f32(Truncate(rand.NormFloat64()*f64(std)+f64(mean), t))
}

// The positive half (when X>mean)  of the custom normal dist :
//
// Custom normal dist variable `X ~ N(mean, std)“.
// Truncated to t decimal places.
func PHalfNormFloat32T(mean, std f32, t int) f32 {
	X := NormFloat32T(mean, std, t)
	if X < mean {
		// X = X + 2*(mean-X)
		X = 2*mean - X
	}
	return X
}

// The negative half (when X<mean)  of the custom normal dist :
//
// Custom normal dist variable `X ~ N(mean, std)“.
// Truncated to t decimal places.
func NHalfNormFloat32T(mean, std f32, t int) f32 {
	X := NormFloat32T(mean, std, t)
	if X > mean {
		// X = X - 2*(X-mean)
		X = 2*mean - X
	}
	return X
}

func RandChoice[T any](a, b T) T {
	if rand.Float32() > 0.5 {
		return a
	} else {
		return b
	}
}

// Absolute of a f32
func Abs(f f32) f32 {
	if f < 0 {
		return -f
	} else {
		return f
	}
}
