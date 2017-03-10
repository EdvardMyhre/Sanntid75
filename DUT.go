package main

// import "./driver"
 import "./types"
// import "./elevator"

import "fmt"
// import "time"
// import "math/rand"



func main(){
	var a []int
	for i:=0; i<10; i++{
		a = append(a, i)
	}
	var b []int
	b = a
	fmt.Println(a)
	a = nil
	fmt.Println(a)
	fmt.Println(b)

	var assigned_tasks []types.Task
	fmt.Println(assigned_tasks)

}