package amanager

import "../driver"
import "../types"

import "fmt"
import "time"

func AssignedTasksManager(elev_status_c <-chan types.Status, elev_task_c chan<- int,
	pmanager_task_c <-chan types.Task, pmanager_status_c chan<- types.Task,
	udp_rx_c <-chan types.MainData, udp_tx_c chan<- types.MainData) {

	//Initializing variables
	var msg_in types.MainData
	var msg_out types.MainData
	var assigned_tasks []types.Task
	var task_temp []types.Task
	var data_temp [][]int
	var task_new types.Task
	var task_current types.Task

	elev_status := types.Status{Destination_floor: 0, Floor: 0, Prev_floor: 1, Finished: 1, Between_floors: 0}
	weight := 255

	//Boot routine
	for button_type := 0; button_type < 3; button_type++ {
		for floor := 0; floor < types.NUMBER_OF_FLOORS; floor++ {
			driver.SetButtonLamp(button_type, floor, 0)
		}
	}

	msg_out.Destination = "backup"
	msg_out.Type = types.REQUEST_BACKUP

	recieve_tries := 10
	send_tries := 10
	loaded := 0
	fmt.Println("AMANAGER: Sending request for backup")

	for j := 0; j < send_tries; j++ {
		if loaded == 1 {
			break
		}

		udp_tx_c <- msg_out
		for i := 0; i < recieve_tries; i++ {
			if loaded == 1 {
				break
			}
			select {
			case msg_in = <-udp_rx_c:
				if msg_in.Type == types.GIVE_BACKUP {
					assigned_tasks = slice2tasks(msg_in.Data)
					loaded = 1
				}
			case <-time.After(time.Second):
				if i == recieve_tries-1 {
					fmt.Println("AMANAGER: Resending request for backup")
				}
			}
		}
	}
	if loaded > 0 {
		fmt.Println("AMANAGER: Backup loaded")
	} else {
		fmt.Println("AMANAGER: No backup! TASKS ARE LOST")
	}
	l := len(assigned_tasks)
	for i := 0; i < l; i++ {
		driver.SetButtonLamp(assigned_tasks[i].Type, assigned_tasks[i].Floor, 1)
	}

	//Reading channels for input
	for {
		//Input from controller, i.e. new status from controller
		select {
		case elev_status = <-elev_status_c:
			if elev_status.Finished != 0 && elev_status.Destination_floor == task_current.Floor {

				//Update assigned tasks
				task_current.Finished = 255
				assigned_tasks = assigned_tasks[1:]

				//Push backup
				msg_out = types.MainData{Destination: "backup", Type: types.PUSH_BACKUP, Data: tasks2slice(assigned_tasks)}
				select {
				case udp_tx_c <- msg_out:
				case <-time.After(time.Second):
					fmt.Println("AMANAGER: network not responding!")
				}

				//Update lights
				if task_current.Type != types.BTN_TYPE_COMMAND {
					msg_out = types.MainData{Destination: "broadcast", Type: types.SET_LIGHT, Data: tasks2slice([]types.Task{task_current})}
					select {
					case udp_tx_c <- msg_out:
					case <-time.After(time.Second):
						fmt.Println("AMANAGER: network not responding!")
					}
				}
				driver.SetButtonLamp(task_current.Type, task_current.Floor, 0)

			}
		case <-time.After(time.Millisecond * 10):
		}

		//Start on new task of we finished the last
		if task_current.Finished != 0 {
			if len(assigned_tasks) > 0 {
				fmt.Println("AMANAGER: starting on new task")
				task_current = assigned_tasks[0]
				select {
				case elev_task_c <- task_current.Floor:
				case <-time.After(time.Second):
					fmt.Println("AMANAGER: elevator.Controller not responding! TASKS ARE LOST")
				}
			}
		}

		//Input from pmanager, i.e new task from pmanager, command from cab or timed out tasks!
		select {
		case task_new = <-pmanager_task_c:
			if task_new.Assigned != 0 {
				fmt.Println("AMANAGER: been assigned an already assigned task!")
			}
			task_new.Assigned = 255
			_, assigned_tasks = addTask(assigned_tasks, task_new, elev_status)

			//Push backup
			msg_out = types.MainData{Destination: "backup", Type: types.PUSH_BACKUP, Data: tasks2slice(assigned_tasks)}
			select {
			case udp_tx_c <- msg_out:
			case <-time.After(time.Second):
				fmt.Println("AMANAGER: network not responding!")
			}

			//Update lights
			if task_new.Type != types.BTN_TYPE_COMMAND {
				msg_out = types.MainData{Destination: "broadcast", Type: types.SET_LIGHT, Data: tasks2slice([]types.Task{task_new})}
				select {
				case udp_tx_c <- msg_out:
				case <-time.After(time.Second):
					fmt.Println("AMANAGER: network not responding!")
				}
			}
			driver.SetButtonLamp(task_new.Type, task_new.Floor, 1)

			//Reply to pmanager
			select {
			case pmanager_status_c <- task_new:
			case <-time.After(time.Second):
				fmt.Println("AMANAGER: pmanager not responding!")
			}
		case <-time.After(time.Millisecond * 10):
		}

		//Message from udp
		select {
		case msg_in = <-udp_rx_c:
			switch msg_in.Type {

			//Weight request
			case types.REQUEST_WEIGHT:
				task_temp = nil
				task_temp = slice2tasks(msg_in.Data)
				if len(task_temp) > 0 {
					task_new = task_temp[0]

					//Send to pmanager
					if task_new.Assigned != 0 {
						fmt.Println("AMANAGER: weight request on assigned task!")
					}
					select {
					case pmanager_status_c <- task_new:
					case <-time.After(time.Second):
						fmt.Println("AMANAGER: pmanager not responding!")
					}

					//Calculate weight
					weight, _ = addTask(assigned_tasks, task_new, elev_status)

					//Send weight to source
					msg_out = types.MainData{Destination: msg_in.Source, Type: types.GIVE_WEIGHT}
					data_temp = make([][]int, 1)
					data_temp[0] = make([]int, 3)
					data_temp[0][0] = weight
					data_temp[0][1] = task_new.Type
					data_temp[0][2] = task_new.Floor
					msg_out.Data = data_temp
					select {
					case udp_tx_c <- msg_out:
					case <-time.After(time.Second):
						fmt.Println("AMANAGER: network not responding!")
					}

				} else {
					fmt.Println("AMANAGER: could not deserialize task from udp!")
				}

			//New task assigned to this elevator
			case types.DISTRIBUTE_ORDER:
				task_temp = nil
				task_temp = slice2tasks(msg_in.Data)
				if len(task_temp) > 0 {
					task_new = task_temp[0]
					if task_new.Assigned != 0 {
						fmt.Println("AMANAGER: been assigned an already assigned task!")
					}
					if task_new.Finished != 0 {
						fmt.Println("AMANAGER: been assigned an already finished task!")
					}
					task_new.Assigned = 255
					_, assigned_tasks = addTask(assigned_tasks, task_new, elev_status)

					//Push backup
					msg_out = types.MainData{Destination: "backup", Type: types.PUSH_BACKUP, Data: tasks2slice(assigned_tasks)}
					select {
					case udp_tx_c <- msg_out:
					case <-time.After(time.Second):
						fmt.Println("AMANAGER: network not responding!")
					}

					//Send to pmanager
					select {
					case pmanager_status_c <- task_new:
					case <-time.After(time.Second):
						fmt.Println("AMANAGER: pmanager not responding!")
					}

					//Confirm that the task is taken
					msg_out = types.MainData{Destination: "broadcast", Type: types.TASK_ASSIGNED, Data: tasks2slice([]types.Task{task_new})}
					select {
					case udp_tx_c <- msg_out:
					case <-time.After(time.Second):
						fmt.Println("AMANAGER: network not responding!")
					}

					//Set local lights
					driver.SetButtonLamp(task_new.Type, task_new.Floor, 1)

				} else {
					fmt.Println("AMANAGER: could not deserialize task from udp!")
				}

			//Asked to set lights
			case types.SET_LIGHT:
				task_temp = nil
				task_temp = slice2tasks(msg_in.Data)
				if len(task_temp) > 0 {
					task_new = task_temp[0]
					if task_new.Assigned == 0 {
						fmt.Println("AMANAGER: has been told to turn on light for unassigned task!")
					}
					if task_new.Type != types.BTN_TYPE_COMMAND {
						if task_new.Finished != 0 {
							driver.SetButtonLamp(task_new.Type, task_new.Floor, 0)
						} else {
							driver.SetButtonLamp(task_new.Type, task_new.Floor, 1)
						}

					}
				} else {
					fmt.Println("AMANAGER: could not deserialize task from udp!")
				}

			//A task has been assigned to another elevator, we inform pmanager and set lights
			case types.TASK_ASSIGNED:
				task_temp = nil
				task_temp = slice2tasks(msg_in.Data)
				if len(task_temp) > 0 {
					task_new = task_temp[0]
					if task_new.Assigned == 0 {
						fmt.Println("AMANAGER: was told an unassigned task has been assigned!")
					}

					//Inform pmanager
					select {
					case pmanager_status_c <- task_new:
					case <-time.After(time.Second):
						fmt.Println("AMANAGER: pmanager not responding!")
					}

					//Update lights
					if task_new.Type != types.BTN_TYPE_COMMAND {
						if task_new.Finished != 0 {
							driver.SetButtonLamp(task_new.Type, task_new.Floor, 0)
						} else {
							driver.SetButtonLamp(task_new.Type, task_new.Floor, 1)
						}

					}

				} else {
					fmt.Println("AMANAGER: could not deserialize task from udp!")
				}

			default:
				fmt.Println("AMANAGER: message from udp unreconisible!")
			}

		case <-time.After(time.Millisecond * 10):
		}
	} //end of inf loop
}

//Functions
func slice2tasks(slice [][]int) []types.Task {
	l := len(slice)
	tasks := make([]types.Task, l)
	for i := 0; i < len(slice); i++ {
		tasks[i].Type = slice[i][0]
		tasks[i].Floor = slice[i][1]
		tasks[i].Finished = slice[i][2]
		tasks[i].Assigned = slice[i][3]
	}
	return tasks
}

func tasks2slice(tasks []types.Task) [][]int {
	l := len(tasks)
	w := 4
	slice := make([][]int, l)
	for i := 0; i < l; i++ {
		slice[i] = make([]int, w)
	}

	for i := 0; i < len(tasks); i++ {
		slice[i][0] = tasks[i].Type
		slice[i][1] = tasks[i].Floor
		slice[i][2] = tasks[i].Finished
		slice[i][3] = tasks[i].Assigned
	}
	return slice
}

func addTask(tasks []types.Task, task types.Task, status types.Status) (int, []types.Task) {

	distance := taskDistance(task, status)
	added := 0
	queuePos := 0
	weight := 0
	posWeight := 0

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
	weight = distance + (queuePos-1)*posWeight

	return weight, tasks
}

func taskDistance(task types.Task, status types.Status) int {
	distance_floors := 0
	abs_distance_floors := 0
	distance_travel := 0

	//Finner bevegelsesretning
	direction := UP
	if status.Prev_floor > status.Floor || status.Floor == types.NUMBER_OF_FLOORS-1 {
		direction = DOWN
	}
	if status.Floor == 0 {
		direction = UP
	}

	//The weight is calculated as if the elevator already is on the destination floor
	if status.Destination_floor == types.NUMBER_OF_FLOORS-1 {
		direction = DOWN
	} else if status.Destination_floor == 0 {
		direction = UP
	}

	distance_floors = task.Floor - status.Destination_floor
	abs_distance_floors = distance_floors
	if distance_floors < 0 {
		abs_distance_floors = -distance_floors
	}

	if direction == UP && distance_floors < 0 {
		distance_travel = 2*(types.NUMBER_OF_FLOORS-1-status.Destination_floor) + abs_distance_floors
	} else if direction == DOWN && distance_floors > 0 {
		distance_travel = 2*status.Destination_floor + abs_distance_floors
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
