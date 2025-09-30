package main

import (
	"fmt"
	"testing"
)

type f32 = float32
type f64 = float64
type u8 = uint8
type u16 = uint16
type u32 = uint32
type u64 = uint64
type i8 = int8
type i16 = int16
type i32 = int32
type i64 = int64

type Pair struct {
	a, b any
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

func main() {
	ob := OrderBook{}
	fmt.Println(ob.GetQueue(ASK))
}
