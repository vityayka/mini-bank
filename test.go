package main

import (
	"fmt"
	"math/rand"
	"time"
)

func test() {
	c1 := chat("Azaz")
	c2 := chat("Kek")
	// timeout := time.After(5 * time.Second)
	for {
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("Timeout, exit")
			return
		case msg := <-c1:
			fmt.Println(msg)
		case msg := <-c2:
			fmt.Println(msg)
		}
	}

}

func chat(name string) <-chan string {
	chat := make(chan string)
	go func() {
		for {
			timeout := rand.Int63n(7000)
			time.Sleep(time.Duration(timeout) * time.Millisecond)
			chat <- fmt.Sprintf("It took %d milliseconds to %s", timeout, name)
		}
	}()
	return chat
}
