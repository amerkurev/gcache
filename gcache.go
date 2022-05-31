package gcache

import (
	"fmt"
)

// PrintSlice prints any slice in generic manner
func PrintSlice[T any](s []T) {
	for _, v := range s {
		fmt.Println(v)
	}
}
