package main

//import "./driver"
//import "./types"

//import "./elevator"
//import "./amanager"

import "fmt"

//import "time"
//import "math/rand"

func main() {
	x := 0
	a := []int{9, 8, 7}
	x, a = a[0], a[1:]
	fmt.Println(x)
	fmt.Println(a)
}
