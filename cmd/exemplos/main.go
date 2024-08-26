package main

import (
	"fmt"
	"time"
)

func worker(workerId int, data chan int) {
	for x:=range(data) {
		fmt.Printf("worker %d got %d \n", workerId, x)
		time.Sleep(time.Second)
	}
}

func main() { //go routine 1
	channel := make(chan int)
	qtdWorkers := 10

	for i:= range(qtdWorkers) {
		go worker(i, channel)
	}
	
	for i:= range(15) {
		channel <- i
	}
}
