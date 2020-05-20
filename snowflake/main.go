package main

import (
	"fmt"
	"time"
)

func main() {
	generate := Constructor(1, 1)
	for i := 0; i < 100; i++ {
		go func() {
			for {
				fmt.Println(generate.nextId())
			}
		}()
	}

	time.Sleep(10 * time.Second)
}
