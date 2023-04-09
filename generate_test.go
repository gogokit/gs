package gs

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFmt(t *testing.T) {
	Convey("AddGetAndSet", t, func() {
		Convey("input_no_lock_filed", func() {
			const src = `
package main
type CommonStruct struct {
	IntSlice1    []int
	IntSlice2    []int
	IntSlice3    []int
	StructPtr    ***CommonStruct
	StructSlice2 []CommonStruct
	StructSlice3 []CommonStruct
	intPtr       ***int
	Map          map[int][]CommonStruct
}
`
			const loopCnt = 1000
			const expectStr = `package main

import "sync"

type CommonStruct struct {
	IntSlice1    []int
	IntSlice2    []int
	IntSlice3    []int
	StructPtr    ***CommonStruct
	StructSlice2 []CommonStruct
	StructSlice3 []CommonStruct
	intPtr       ***int
	Map          map[int][]CommonStruct
	lock         sync.RWMutex
}

func (obj *CommonStruct) GetIntSlice1() []int {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.IntSlice1
}
func (obj *CommonStruct) SetIntSlice1(param []int) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.IntSlice1 = param
	return obj
}
func (obj *CommonStruct) GetIntSlice2() []int {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.IntSlice2
}
func (obj *CommonStruct) SetIntSlice2(param []int) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.IntSlice2 = param
	return obj
}
func (obj *CommonStruct) GetIntSlice3() []int {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.IntSlice3
}
func (obj *CommonStruct) SetIntSlice3(param []int) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.IntSlice3 = param
	return obj
}
func (obj *CommonStruct) GetStructPtr() ***CommonStruct {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.StructPtr
}
func (obj *CommonStruct) SetStructPtr(param ***CommonStruct) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.StructPtr = param
	return obj
}
func (obj *CommonStruct) GetStructSlice2() []CommonStruct {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.StructSlice2
}
func (obj *CommonStruct) SetStructSlice2(param []CommonStruct) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.StructSlice2 = param
	return obj
}
func (obj *CommonStruct) GetStructSlice3() []CommonStruct {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.StructSlice3
}
func (obj *CommonStruct) SetStructSlice3(param []CommonStruct) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.StructSlice3 = param
	return obj
}
func (obj *CommonStruct) GetIntPtr() ***int {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.intPtr
}
func (obj *CommonStruct) SetIntPtr(param ***int) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.intPtr = param
	return obj
}
func (obj *CommonStruct) GetMap() map[int][]CommonStruct {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.Map
}
func (obj *CommonStruct) SetMap(param map[int][]CommonStruct) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.Map = param
	return obj
}
`
			for i := 1; i <= loopCnt; i++ {
				str, err := AddGetAndSet(src)
				So(err, ShouldEqual, nil)
				So(str, ShouldEqual, expectStr)
			}
		})

		Convey("input_had_lock_filed", func() {
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
	lock         sync.Mutex
}
`
			const loopCnt = 1000
			const expectStr = `package main

type CommonStruct struct {
	IntSlice1    []int
	IntSlice2    []int
	IntSlice3    []int
	StructPtr    ***CommonStruct
	StructSlice2 []CommonStruct
	StructSlice3 []CommonStruct
	intPtr       ***int
	Map          map[int][]CommonStruct
	lock         sync.Mutex
}
`
			for i := 1; i <= loopCnt; i++ {
				str, err := AddGetAndSet(src)
				So(err, ShouldEqual, nil)
				So(str, ShouldEqual, expectStr)
			}
		})

		Convey("input_had_some_method_filed", func() {
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
}

func (obj *CommonStruct) GetIntSlice1() {
	return
}
`
			const loopCnt = 1000
			const expectStr = `package main

import "sync"

type CommonStruct struct {
	IntSlice1    []int
	IntSlice2    []int
	IntSlice3    []int
	StructPtr    ***CommonStruct
	StructSlice2 []CommonStruct
	StructSlice3 []CommonStruct
	intPtr       ***int
	Map          map[int][]CommonStruct
	lock         sync.RWMutex
}

func (obj *CommonStruct) GetIntSlice1() {
	return
}
func (obj *CommonStruct) SetIntSlice1(param []int) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.IntSlice1 = param
	return obj
}
func (obj *CommonStruct) GetIntSlice2() []int {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.IntSlice2
}
func (obj *CommonStruct) SetIntSlice2(param []int) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.IntSlice2 = param
	return obj
}
func (obj *CommonStruct) GetIntSlice3() []int {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.IntSlice3
}
func (obj *CommonStruct) SetIntSlice3(param []int) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.IntSlice3 = param
	return obj
}
func (obj *CommonStruct) GetStructPtr() ***CommonStruct {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.StructPtr
}
func (obj *CommonStruct) SetStructPtr(param ***CommonStruct) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.StructPtr = param
	return obj
}
func (obj *CommonStruct) GetStructSlice2() []CommonStruct {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.StructSlice2
}
func (obj *CommonStruct) SetStructSlice2(param []CommonStruct) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.StructSlice2 = param
	return obj
}
func (obj *CommonStruct) GetStructSlice3() []CommonStruct {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.StructSlice3
}
func (obj *CommonStruct) SetStructSlice3(param []CommonStruct) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.StructSlice3 = param
	return obj
}
func (obj *CommonStruct) GetIntPtr() ***int {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.intPtr
}
func (obj *CommonStruct) SetIntPtr(param ***int) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.intPtr = param
	return obj
}
func (obj *CommonStruct) GetMap() map[int][]CommonStruct {
	obj.lock.RLock()
	defer obj.lock.RUnlock()
	return obj.Map
}
func (obj *CommonStruct) SetMap(param map[int][]CommonStruct) *CommonStruct {
	if obj == nil {
		return nil
	}
	obj.lock.Lock()
	defer obj.lock.Unlock()
	obj.Map = param
	return obj
}
`
			for i := 1; i <= loopCnt; i++ {
				str, err := AddGetAndSet(src)
				So(err, ShouldEqual, nil)
				So(str, ShouldEqual, expectStr)
			}
		})
	})
}
