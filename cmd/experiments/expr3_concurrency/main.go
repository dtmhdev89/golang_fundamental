package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Go concurrency example")
	go execute("Routine 1")
	go execute("Routine 2")
	go execute("Routine 3")

	varChannel := make(chan int)
	go sharetoChannel(varChannel)
	go getFromChannel(varChannel)

	time.Sleep(time.Second)
	fmt.Println("Example completed")
}

func execute(message string) {
	for i := 1; i <= 5; i++ {
		time.Sleep(50 * time.Microsecond)
		fmt.Println(message, ": ", i)
	}
}

func sharetoChannel(varChannel chan int) {
	for i := 0; i <= 7; i++ {
		varChannel <- i
	}
	close(varChannel)
}

func getFromChannel(varChannel chan int) {
	for i := 0; i <= 7; i++ {
		fmt.Println(<-varChannel)
	}
}
