package main

//import "./driver"
import "./types"

//import "./elevator"
//import "./amanager"

import "fmt"

//import "time"
//import "math/rand"

func main() {
	var a [][]int
	a = make([][]int, 1)
	a[0] = make([]int, 3)

	var msg_in types.MainData
	var msg_a types.MainData
	msg_a = types.MainData{Destination: "backup", Data: a}
	msg_in = types.MainData{Destination: "lol", Data: msg_a.Data}
	fmt.Println(msg_in)
}
