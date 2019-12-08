package main

import (
	"fmt"
	"time"
)

func printString(c chan string, no int) {

	for v := range c {
		fmt.Println(no, v)
		time.Sleep(time.Millisecond)
	}
}

func main() {
	ch := make(chan string, 10)
	for i := 0; i < 10; i++ {
		go printString(ch, i)
	}

	time.Sleep(time.Second)
	for i := 0; i < 50; i++ {
		ch <- "hello"
	}
	time.Sleep(time.Second)
}
