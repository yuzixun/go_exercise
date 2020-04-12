package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

type Foo struct {
	a int
}

func NewFoo(i int) *Foo {
	f := &Foo{a: rand.Intn(50)}
	runtime.SetFinalizer(f, func(f *Foo) {
		fmt.Println(`foo ` + strconv.Itoa(i) + ` has been gc`)
	})

	return f
}

func main() {
	for i := 0; i < 3; i++ {
		f := NewFoo(i)
		println(f.a)
	}

	runtime.GC()
	// 无法保证在程序退出前执行，因此添加等待
	time.Sleep(time.Second)
}
