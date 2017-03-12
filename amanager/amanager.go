package amanager

//import "../driver"
//import "../elevator"
import "../types"

import "fmt"
import "time"

func AssignedTasksManager(elev_status_c <-chan types.Status, elev_task_c chan<- int, statusc <-chan types.Task, taskc chan<- types.Task, udp_rx_c <-chan types.MainData, udp_tx_c chan<- types.MainData) {

	//Initializing variables
	var assigned_tasks []types.Task
	var msg types.MainData

	//Boot routine
	msg.Destination = "backup"
	msg.Message_type = types.MESSAGE_TYPE_REQUEST_BACKUP
	udp_tx_c <- msg
	select {
	case msg := <-udp_rx_c:
		assigned_tasks = slice2tasks(msg.Data)
		fmt.Println("AMANAGER: Backup loaded")
	case <-time.After(time.Second * 10):
		fmt.Println("AMANAGER: No backup available")
	}

	//HENDELSER:
	//Får task fra pmanager
	//Får task fra udp
	//Får spm om vekt fra udp
	//Får status fra controller
	//

}

//Functions
func slice2tasks(slice [][]int) []types.Task {
	l := len(slice)
	tasks := make([]types.Task, l)
	for i := 0; i < len(slice); i++ {
		tasks[i].Type = slice[i][0]
		tasks[i].Floor = slice[i][1]
		tasks[i].Add = slice[i][2]
	}
	return tasks
}

func tasks2slice(tasks []types.Task) [][]int {
	l := len(tasks)
	w := 3
	slice := make([][]int, l)
	for i := 0; i < l; i++ {
		slice[i] = make([]int, w)
	}

	for i := 0; i < len(tasks); i++ {
		slice[i][0] = tasks[i].Type
		slice[i][1] = tasks[i].Floor
		slice[i][2] = tasks[i].Add
	}
	return slice
}

func addTask(tasks []types.Task, task types.Task, status types.Status) (int, []types.Task) {
	//task ligger ikke i tasks fra før av. Skal nå legge den inn i riktig rekkefølge og returnere en vekt på å utføre denne task.
	//Ingen merging av tasks

	distance := taskDistance(task, status)
	added := 0
	queuePos := 0
	weight := 0
	posWeight := 3

	l := len(tasks)
	for i := 0; i < l; i++ {
		if distance < taskDistance(tasks[i], status) {
			tasks = append(tasks[:i], append([]types.Task{task}, tasks[i:]...)...)
			added = 1
			queuePos = i + 1
			break
		}
	}
	if added == 0 {
		tasks = append(tasks, task)
		queuePos = len(tasks)
	}

	fmt.Println("queuePos:", queuePos)
	fmt.Println("distance:", distance)
	weight = distance + (queuePos-1)*posWeight

	return weight, tasks
}

func taskDistance(task types.Task, status types.Status) int {
	distance_floors := 0
	abs_distance_floors := 0
	distance_travel := 0
	floor_offset := status.Floor

	//Finner bevegelsesretning
	direction := UP
	if status.Prev_floor > status.Floor || status.Floor == types.NUMBER_OF_FLOORS-1 {
		direction = DOWN
	}
	if status.Floor == 0 {
		direction = UP
	}

	//The weight is calculated as if the elevator already is on the first floor in its direction of travel
	if status.Between_floors == 1 {
		switch direction {
		case UP:
			floor_offset += 1
		case DOWN:
			floor_offset -= 1
		}
	}
	if floor_offset == types.NUMBER_OF_FLOORS-1 {
		direction = DOWN
	} else if floor_offset == 0 {
		direction = UP
	}

	if floor_offset < 0 || floor_offset > types.NUMBER_OF_FLOORS-1 {
		fmt.Println("ERROR: Floor out of bounds!")
	}

	distance_floors = task.Floor - floor_offset
	abs_distance_floors = distance_floors
	if distance_floors < 0 {
		abs_distance_floors = -distance_floors
	}

	if direction == UP && distance_floors < 0 {
		distance_travel = 2*(types.NUMBER_OF_FLOORS-1-floor_offset) + abs_distance_floors
	} else if direction == DOWN && distance_floors > 0 {
		distance_travel = 2*floor_offset + abs_distance_floors
	} else {
		distance_travel = abs_distance_floors
	}

	if direction == UP && task.Type == 0 && distance_floors < 0 {
		distance_travel += 2 * task.Floor
	} else if direction == UP && task.Type == 1 && distance_floors >= 0 {
		distance_travel += 2 * (types.NUMBER_OF_FLOORS - 1 - task.Floor)
	} else if direction == DOWN && task.Type == 1 && distance_floors > 0 {
		distance_travel += 2 * (types.NUMBER_OF_FLOORS - 1 - task.Floor)
	} else if direction == DOWN && task.Type == 0 && distance_floors <= 0 {
		distance_travel += 2 * task.Floor
	}

	if distance_travel > 2*types.NUMBER_OF_FLOORS-3 {
		fmt.Println("ERROR: Distance is longer than physically possible!")
	}
	return distance_travel
}

const (
	UP   = 1
	DOWN = 0
)

//TEST QUEUING
// func main() {
// 	var tasks []types.Task
// 	weight := 0

// 	status := types.Status{Destination_floor: 3, Floor: 1, Prev_floor: 0, Finished: 0, Between_floors: 1}

// 	task1 := types.Task{Type: 0, Floor: 0, Add: 255}
// 	task2 := types.Task{Type: 0, Floor: 1, Add: 255}
// 	task3 := types.Task{Type: 0, Floor: 2, Add: 255}
// 	task4 := types.Task{Type: 1, Floor: 3, Add: 255}
// 	task5 := types.Task{Type: 1, Floor: 2, Add: 255}
// 	task6 := types.Task{Type: 1, Floor: 1, Add: 255}

// 	weight, tasks = amanager.AddTask(tasks, task1, status)
// 	fmt.Println("Weight:", weight)
// 	weight, tasks = amanager.AddTask(tasks, task2, status)
// 	fmt.Println("Weight:", weight)
// 	weight, tasks = amanager.AddTask(tasks, task3, status)
// 	fmt.Println("Weight:", weight)
// 	weight, tasks = amanager.AddTask(tasks, task4, status)
// 	fmt.Println("Weight:", weight)
// 	weight, tasks = amanager.AddTask(tasks, task5, status)
// 	fmt.Println("Weight:", weight)
// 	weight, tasks = amanager.AddTask(tasks, task6, status)

// 	fmt.Println("Tasks:")
// 	for i := 0; i < len(tasks); i++ {
// 		fmt.Println(tasks[i])
// 	}
// 	fmt.Println("Weight:", weight)
// }
