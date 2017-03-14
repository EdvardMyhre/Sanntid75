package main

import (
	//"fmt"
	"time"
)

import (
	"./elevator"
	"./types"
)

func main() {
	chan_button := make(chan types.Task)
	go elevator.ButtonPoller(chan_button)
	for {
		time.Sleep(time.Second)
	}
}
