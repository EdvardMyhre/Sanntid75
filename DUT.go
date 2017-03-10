package main

//import "./driver"
// import "./types"
//import "./elevator"

import "fmt"
//import "time"
//import "math/rand"



//TEST CONTROLLER
func main() {
	var a [10][10]

	for i := 0; i < 10; i++{
		a[i] = append(a[i], 0)
		a[i] = append(a[i], 1)
		a[i] = append(a[i], 2)
	}

	fmt.Println(len(a))
	fmt.Println(lengde(a))

	// var a [][]int
	// var b []types.Task
	// c := types.Task{Type: 1, Floor: 1, Add: 255}
	// b = append(b, c)
	// b = append(b, c)
	// fmt.Println(b)
	// fmt.Println(len(b))
	// //a = tasks2slice(b)
	// //fmt.Println(a)
}

func lengde(a [][]int) int{
	return len(a)
}

// func tasks2slice(tasks []types.Task) [][]int{
// 	var slice [len(tasks)][3]int
// 	for i := 0; i < len(tasks); i++{
// 		slice[i][0] = tasks[i].Type
// 		slice[i][1] = tasks[i].Floor
// 		slice[i][2] = tasks[i].Add
// 	}
// 	return slice
// }