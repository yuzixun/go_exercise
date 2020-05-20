package main

import (
	"fmt"
	"time"
)

func main() {
	idChan := make(chan int64)
	generate := Constructor(1, 1)
	for i := 0; i < 100; i++ {
		go func() {
			for {
				idChan <- generate.nextId()
			}
		}()
	}

	go func(idChan chan int64) {
		for {
			fmt.Println(<-idChan)
		}
	}(idChan)

	time.Sleep(10 * time.Second)
}
