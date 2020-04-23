package main

import (
	"fmt"
	"sync"
)

var pool *sync.Pool

type Person struct {
	Name string
}

func initPool() {
	pool = &sync.Pool{
		New: func() interface{} {
			fmt.Println("create a new person")
			return new(Person)
		},
	}
}

func main() {
	initPool()

	p := pool.Get().(*Person)
	fmt.Println("get from pool:", p)

	p.Name = "first"
	fmt.Println("set p.name = ", p.Name)

	pool.Put(p)

	fmt.Println("pool has an object, get it", pool.Get().(*Person))
	fmt.Println("pool has no object, get", pool.Get().(*Person))
}
