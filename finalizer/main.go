package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

type Foo struct {
	a int
}

func NewFoo(i int) *Foo {
	f := &Foo{a: rand.Intn(50)}
	// 在程序无法获取到一个 obj 所指向的对象后的任意时刻，
	// finalizer 被调度运行，且无法保证 finalizer 运行在程序退出之前。
	// 因此一般情况下，因此它们仅用于在长时间运行的程序上释放一些与对象关联的非内存资源。
	runtime.SetFinalizer(f, func(f *Foo) {
		_ = fmt.Sprintf(`foo ` + strconv.Itoa(i) + ` has been gc`)
	})

	return f
}

func main() {
	debug.SetGCPercent(-1)
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	// 初始状态
	fmt.Printf("Allocation: %f Mb, Number of allocation: %d \n", float32(ms.HeapAlloc)/float32(1024*1024), ms.HeapObjects)

	for i := 0; i < 1000000; i++ {
		f := NewFoo(i)
		_ = fmt.Sprintf("%d", f.a)
	}

	// 创建了对象之后的状态
	runtime.ReadMemStats(&ms)
	fmt.Printf("Allocation: %f Mb, Number of allocation: %d \n", float32(ms.HeapAlloc)/float32(1024*1024), ms.HeapObjects)

	runtime.GC()
	time.Sleep(time.Second)

	// 第一次GC完成，释放了NewFoo对象。
	// 但是由于SetFinalizer的存在，
	// 创建了新的GoRoutine，并且使得NewFoo对象再次可访问
	runtime.ReadMemStats(&ms)
	fmt.Printf("Allocation: %f Mb, Number of allocation: %d \n", float32(ms.HeapAlloc)/float32(1024*1024), ms.HeapObjects)

	runtime.GC()
	time.Sleep(time.Second)

	// 最后一次GC清除了NewFoo
	runtime.ReadMemStats(&ms)
	fmt.Printf("Allocation: %f Mb, Number of allocation: %d \n", float32(ms.HeapAlloc)/float32(1024*1024), ms.HeapObjects)

	// runtime.GC()
	// time.Sleep(time.Second)

}
