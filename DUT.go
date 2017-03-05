package main

import "./driver"
import "./types"
import "./elevator"

import "fmt"
//import "time"
//import "math/rand"


func main() {
	driver.Init()
	buttonc := make(chan types.Button)
	go elevator.ButtonPoller(buttonc)
	for{
		button := <- buttonc
		fmt.Println("Type: ", button.Type, " Floor: ", button.Floor)
	}

}
