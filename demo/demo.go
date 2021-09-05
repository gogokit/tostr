package main

import (
	"fmt"

	"github.com/gogokit/tostr"
)

func main() {
	type Struct struct {
		slice1 []int
		slice2 []int
	}
	slice := []int{1, 2, 3, 4, 5}
	s := Struct{
		slice1: slice,
		slice2: slice[1:3],
	}
	fmt.Printf("%v\n", tostr.String(s))
	fs, _ := tostr.Fmt(tostr.String(s), 2)
	fmt.Printf("\n%v\n", fs)
}
