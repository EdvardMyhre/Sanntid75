package main

//import "./driver"
import "./types"

//import "./elevator"
import "./amanager"

import "fmt"

//import "time"
//import "math/rand"

func main() {
	var tasks []types.Task
	weight := 0

	status := types.Status{Destination_floor: 3, Floor: 1, Prev_floor: 0, Finished: 0, Between_floors: 1}

	task1 := types.Task{Type: 0, Floor: 0, Add: 255}
	task2 := types.Task{Type: 0, Floor: 1, Add: 255}
	task3 := types.Task{Type: 0, Floor: 2, Add: 255}
	task4 := types.Task{Type: 1, Floor: 3, Add: 255}
	task5 := types.Task{Type: 1, Floor: 2, Add: 255}
	task6 := types.Task{Type: 1, Floor: 1, Add: 255}

	weight, tasks = amanager.AddTask(tasks, task1, status)
	fmt.Println("Weight:", weight)
	weight, _ = amanager.AddTask(tasks, task2, status)
	fmt.Println("Weight:", weight)
	weight, _ = amanager.AddTask(tasks, task3, status)
	fmt.Println("Weight:", weight)
	weight, tasks = amanager.AddTask(tasks, task4, status)
	fmt.Println("Weight:", weight)
	weight, _ = amanager.AddTask(tasks, task5, status)
	fmt.Println("Weight:", weight)
	weight, _ = amanager.AddTask(tasks, task6, status)

	fmt.Println("Tasks:")
	for i := 0; i < len(tasks); i++ {
		fmt.Println(tasks[i])
	}
	fmt.Println("Weight:", weight)

	deleteTasks(tasks)
	fmt.Println(tasks)
}
