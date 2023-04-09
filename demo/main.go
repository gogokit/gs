package main

import (
	"fmt"

	"github.com/gogokit/gs"
)

const src = `package main

type CommonStruct struct {
	IntSlice1    []int
	IntSlice2    []int
	IntSlice3    []int
	StructPtr    ***CommonStruct
	StructSlice2 []CommonStruct
	StructSlice3 []CommonStruct
	intPtr       ***int
	Map          map[int][]CommonStruct
}`

func main() {
	newSrc, err := gs.AddGetAndSet(src)
	if err != nil {
		panic(err)
	}
	fmt.Printf("输入:\n%s\n输出:\n%s\n", src, newSrc)
}
