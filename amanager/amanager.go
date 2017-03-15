package amanager

import "../driver"
import "../types"

import "fmt"
import "time"

func Assigned_tasks_manager(elev_status_c <-chan types.Status, elev_task_c chan<- int,
	pmanager_task_c <-chan types.Task, pmanager_status_c chan<- types.Task,
	udp_rx_c <-chan types.MainData, udp_tx_c chan<- types.MainData,
	chan_backupRecieve <-chan types.MainData) {

	var msg_in types.MainData
	var msg_out types.MainData
	var assigned_tasks []types.Task
	var tasks_temp []types.Task
	var backup_returned []types.Task
	var copy_assigned []types.Task
	var data_temp [][]int
	var task_new types.Task
	var task_current types.Task
	var time_start time.Time
	var alike int
	var index int

	button_lights := make([][]time.Time, types.NUMBER_OF_FLOORS)
	for i := 0; i < types.NUMBER_OF_FLOORS; i++ {
		button_lights[i] = make([]time.Time, 2)
	}

	elev_status := types.Status{Destination_floor: 0, Floor: 0, Prev_floor: 1, Finished: 1, Between_floors: 0}
	weight := 255

	//Boot routine
	for button_type := 0; button_type < 3; button_type++ {
		for floor := 0; floor < types.NUMBER_OF_FLOORS; floor++ {
			driver.SetButtonLamp(button_type, floor, 0)
		}
	}

	fmt.Println("AMANAGER: requesting backup")
	msg_out = types.MainData{Destination: "backup", Type: types.REQUEST_BACKUP}
	time_start = time.Now()
	udp_tx_c <- msg_out
Boot_loop:
	for {
		if time.Since(time_start).Nanoseconds() >= types.TIMEOUT_BACKUP_RESPONSE {
			fmt.Println("AMANAGER: could not reach backup! Tasks might have been lost...")
			break Boot_loop
		}
		select {
		case msg_in = <-chan_backupRecieve:
			fmt.Println("AMAGER: backup response:", msg_in)
			if msg_in.Type == types.GIVE_BACKUP {
				assigned_tasks = slice2tasks(msg_in.Data)
				fmt.Println("AMANGER: backup loaded")
				break Boot_loop
			}
		case <-time.After(types.RETRY_BACKUP_RESPONSE):
			udp_tx_c <- msg_out
		}
	}

	for i := 0; i < len(assigned_tasks); i++ {
		driver.SetButtonLamp(assigned_tasks[i].Type, assigned_tasks[i].Floor, 1)
		if assigned_tasks[i].Type != types.BTN_TYPE_COMMAND {
			button_lights[assigned_tasks[i].Floor][assigned_tasks[i].Type] = time.Now()
		}
	}

	//Reading channels for input
	for {
		//Input from controller
		select {
		case elev_status = <-elev_status_c:
			if elev_status.Finished != 0 {
				if elev_status.Floor != task_current.Floor {
					fmt.Println("AMANAGER: elevator is finished at wrong floor")
				}
				//Update assigned tasks
				task_current.Finished = 255
				tasks_temp = nil
				index = 0
				if len(assigned_tasks) > 0 {
					for {
						if elev_status.Floor == assigned_tasks[index].Floor {
							assigned_tasks[index].Finished = 255
							tasks_temp = append(tasks_temp, assigned_tasks[index])
							assigned_tasks = append(assigned_tasks[:index], assigned_tasks[index+1:]...)
							index--
						}
						if index == len(assigned_tasks)-1 {
							break
						}
						index++
					}
				}

				//Push backup
				msg_out = types.MainData{Destination: "backup", Type: types.PUSH_BACKUP, Data: tasks2slice(assigned_tasks)}
				time_start = time.Now()
				select {
				case udp_tx_c <- msg_out:
				case <-time.After(types.RETRY_BACKUP_RESPONSE):
				}

			Push_backup_1:
				for {
					if time.Since(time_start).Nanoseconds() >= types.TIMEOUT_BACKUP_RESPONSE {
						fmt.Println("AMANAGER: could not reach backup! Tasks might have been lost...")
						break Push_backup_1
					}
					select {
					case msg_in = <-chan_backupRecieve:
						fmt.Println("AMAGER: backup response:", msg_in)
						backup_returned = slice2tasks(msg_in.Data)
						alike = 255
						if len(backup_returned) != len(assigned_tasks) {
							break Push_backup_1
						}
						for i := 0; i < len(backup_returned); i++ {
							if assigned_tasks[i] != backup_returned[i] {
								alike = 0
								break
							}
						}
						if msg_in.Type == types.ACK_BACKUP && alike != 0 {
							fmt.Println("AMANAGER: Backup sucsessfully pushed and ack-ed")
							break Push_backup_1
						}
					case <-time.After(types.RETRY_BACKUP_RESPONSE):
						udp_tx_c <- msg_out
					}
				}

				//Update local lights
				for i := 0; i < len(tasks_temp); i++ {
					driver.SetButtonLamp(tasks_temp[i].Type, tasks_temp[i].Floor, 0)
					if tasks_temp[i].Type != types.BTN_TYPE_COMMAND {
						button_lights[tasks_temp[i].Floor][tasks_temp[i].Type] = time.Since(time.Now())
					}
				}

				//Delete cab commands before updating non-local lights
				index = 0
				if len(tasks_temp) > 0 {
					for {
						if tasks_temp[index].Type == types.BTN_TYPE_COMMAND {
							tasks_temp = append(tasks_temp[:index], tasks_temp[index+1:]...)
							index--
						}
						if index == len(tasks_temp)-1 {
							break
						}
						index++
					}
				}

				msg_out = types.MainData{Destination: "broadcast", Type: types.SET_LIGHT, Data: tasks2slice(tasks_temp)}
				select {
				case udp_tx_c <- msg_out:
				case <-time.After(types.TIMEOUT_AMANAGER_WAITTIME):
					fmt.Println("AMANAGER: network not responding!")
				}

			}
		case <-time.After(types.PAUSE_AMAGER):
		}

		//Start on new task if we are finished with last
		if task_current.Finished != 0 {
			if len(assigned_tasks) > 0 {
				task_current = assigned_tasks[0]
				select {
				case elev_task_c <- task_current.Floor:
				case <-time.After(types.TIMEOUT_AMANAGER_WAITTIME):
					fmt.Println("AMANAGER: elevator.Controller not responding! TASKS ARE LOST")
				}
			}
		}

		//Input from pmanager, i.e new task from pmanager, command from cab or timed-out tasks
		select {
		case task_new = <-pmanager_task_c:
			if task_new.Assigned != 0 {
				fmt.Println("AMANAGER: been assigned an already assigned task!")
			}
			task_new.Assigned = 255

			copy_assigned = nil
			for i := 0; i < len(assigned_tasks); i++ {
				copy_assigned = append(copy_assigned, assigned_tasks[i])
			}
			_, assigned_tasks = addTask(copy_assigned, task_new, elev_status)
			fmt.Println("Assigned tasks after getting task from pmanager")

			//Push backup
			fmt.Println("AMANAGER: pushing backup")
			msg_out = types.MainData{Destination: "backup", Type: types.PUSH_BACKUP, Data: tasks2slice(assigned_tasks)}
			time_start = time.Now()
			select {
			case udp_tx_c <- msg_out:
			case <-time.After(types.RETRY_BACKUP_RESPONSE):
			}
		Push_backup_2:
			for {
				if time.Since(time_start).Nanoseconds() >= types.TIMEOUT_BACKUP_RESPONSE {
					fmt.Println("AMANAGER: could not reach backup! Tasks might have been lost...")
					break Push_backup_2
				}
				select {
				case msg_in = <-chan_backupRecieve:
					fmt.Println("AMAGER: backup response:", msg_in)
					backup_returned = slice2tasks(msg_in.Data)
					alike = 255
					if len(backup_returned) != len(assigned_tasks) {
						break Push_backup_2
					}
					for i := 0; i < len(backup_returned); i++ {
						if assigned_tasks[i] != backup_returned[i] {
							alike = 0
							break
						}
					}
					if msg_in.Type == types.ACK_BACKUP && alike != 0 {
						fmt.Println("AMANAGER: Backup sucsessfully pushed and ack-ed")
						break Push_backup_2
					}
				case <-time.After(types.RETRY_BACKUP_RESPONSE):
					udp_tx_c <- msg_out
				}
			}

			//Update lights
			if task_new.Type != types.BTN_TYPE_COMMAND {
				msg_out = types.MainData{Destination: "broadcast", Type: types.SET_LIGHT, Data: tasks2slice([]types.Task{task_new})}
				select {
				case udp_tx_c <- msg_out:
				case <-time.After(types.TIMEOUT_AMANAGER_WAITTIME):
					fmt.Println("AMANAGER: network not responding!")
				}
			}
			driver.SetButtonLamp(task_new.Type, task_new.Floor, 1)
			if task_new.Type != types.BTN_TYPE_COMMAND {
				button_lights[task_new.Floor][task_new.Type] = time.Now()
			}

			//Reply to pmanager
			select {
			case pmanager_status_c <- task_new:
			case <-time.After(types.TIMEOUT_AMANAGER_WAITTIME):
				fmt.Println("AMANAGER: pmanager not responding!")
			}

		//Message from udp
		case msg_in = <-udp_rx_c:
			switch msg_in.Type {

			//Weight request
			case types.REQUEST_WEIGHT:
				tasks_temp = nil
				tasks_temp = slice2tasks(msg_in.Data)
				if len(tasks_temp) > 0 {
					task_new = tasks_temp[0]

					//Send to pmanager
					if task_new.Assigned != 0 {
						fmt.Println("AMANAGER: weight request on assigned task!")
					}
					select {
					case pmanager_status_c <- task_new:
					case <-time.After(types.TIMEOUT_AMANAGER_WAITTIME):
						fmt.Println("AMANAGER: pmanager not responding!")
					}

					//Calculate weight
					copy_assigned = nil
					for i := 0; i < len(assigned_tasks); i++ {
						copy_assigned = append(copy_assigned, assigned_tasks[i])
					}
					weight, _ = addTask(copy_assigned, task_new, elev_status)

					//Send weight to source
					msg_out = types.MainData{Destination: msg_in.Source, Type: types.GIVE_WEIGHT}
					data_temp = nil
					data_temp = make([][]int, 1)
					data_temp[0] = make([]int, 3)
					data_temp[0][0] = weight
					data_temp[0][1] = task_new.Type
					data_temp[0][2] = task_new.Floor
					msg_out.Data = data_temp
					select {
					case udp_tx_c <- msg_out:
					case <-time.After(types.TIMEOUT_AMANAGER_WAITTIME):
						fmt.Println("AMANAGER: network not responding!")
					}

				} else {
					fmt.Println("AMANAGER: could not deserialize task from udp!")
				}

			//New task assigned to this elevator
			case types.DISTRIBUTE_ORDER:
				tasks_temp = nil
				tasks_temp = slice2tasks(msg_in.Data)
				if len(tasks_temp) > 0 {
					task_new = tasks_temp[0]
					if task_new.Assigned != 0 {
						fmt.Println("AMANAGER: been assigned an already assigned task!")
					}
					if task_new.Finished != 0 {
						fmt.Println("AMANAGER: been assigned an already finished task!")
					}
					task_new.Assigned = 255
					copy_assigned = nil
					for i := 0; i < len(assigned_tasks); i++ {
						copy_assigned = append(copy_assigned, assigned_tasks[i])
					}
					_, assigned_tasks = addTask(copy_assigned, task_new, elev_status)

					//Push backup
					msg_out = types.MainData{Destination: "backup", Type: types.PUSH_BACKUP, Data: tasks2slice(assigned_tasks)}
					time_start = time.Now()
					select {
					case udp_tx_c <- msg_out:
					case <-time.After(types.RETRY_BACKUP_RESPONSE):
					}
				Push_backup_3:
					for {
						if time.Since(time_start).Nanoseconds() >= types.TIMEOUT_BACKUP_RESPONSE {
							fmt.Println("AMANAGER: could not reach backup! Tasks might have been lost...")
							break Push_backup_3
						}
						select {
						case msg_in = <-chan_backupRecieve:
							fmt.Println("AMAGER: backup response:", msg_in)
							backup_returned = slice2tasks(msg_in.Data)
							alike = 255
							if len(backup_returned) != len(assigned_tasks) {
								break Push_backup_3
							}
							for i := 0; i < len(backup_returned); i++ {
								if assigned_tasks[i] != backup_returned[i] {
									alike = 0
									break
								}
							}
							if msg_in.Type == types.ACK_BACKUP && alike != 0 {
								fmt.Println("AMANAGER: Backup sucsessfully pushed and ack-ed")
								break Push_backup_3
							}
						case <-time.After(types.RETRY_BACKUP_RESPONSE):
							udp_tx_c <- msg_out
						}
					}

					//Send to pmanager
					select {
					case pmanager_status_c <- task_new:
					case <-time.After(types.TIMEOUT_AMANAGER_WAITTIME):
						fmt.Println("AMANAGER: pmanager not responding!")
					}

					//Confirm that the task is taken
					msg_out = types.MainData{Destination: "broadcast", Type: types.TASK_ASSIGNED, Data: tasks2slice([]types.Task{task_new})}
					select {
					case udp_tx_c <- msg_out:
					case <-time.After(types.TIMEOUT_AMANAGER_WAITTIME):
						fmt.Println("AMANAGER: network not responding!")
					}

					//Set local lights
					driver.SetButtonLamp(task_new.Type, task_new.Floor, 1)
					if task_new.Type != types.BTN_TYPE_COMMAND {
						button_lights[task_new.Floor][task_new.Type] = time.Now()
					}

				} else {
					fmt.Println("AMANAGER: could not deserialize task from udp!")
				}

			//Asked to set lights
			case types.SET_LIGHT:
				tasks_temp = nil
				tasks_temp = slice2tasks(msg_in.Data)
				if len(tasks_temp) > 0 {
					for i := 0; i < len(tasks_temp); i++ {
						task_new = tasks_temp[i]
						if task_new.Assigned == 0 {
							fmt.Println("AMANAGER: has been told to turn on light for unassigned task!")
						}
						if task_new.Type != types.BTN_TYPE_COMMAND {
							if task_new.Finished != 0 {
								driver.SetButtonLamp(task_new.Type, task_new.Floor, 0)
								button_lights[task_new.Floor][task_new.Type] = time.Since(time.Now())
							} else {
								driver.SetButtonLamp(task_new.Type, task_new.Floor, 1)
								button_lights[task_new.Floor][task_new.Type] = time.Now()
							}

						}
					}
				} else {
					fmt.Println("AMANAGER: could not deserialize task from udp!")
				}

			//A task has been assigned to another elevator
			case types.TASK_ASSIGNED:
				tasks_temp = nil
				tasks_temp = slice2tasks(msg_in.Data)
				if len(tasks_temp) > 0 {
					task_new = tasks_temp[0]
					if task_new.Assigned == 0 {
						fmt.Println("AMANAGER: was told an unassigned task has been assigned!")
					}

					//Inform pmanager
					select {
					case pmanager_status_c <- task_new:
					case <-time.After(types.TIMEOUT_AMANAGER_WAITTIME):
						fmt.Println("AMANAGER: pmanager not responding!")
					}

					//Update lights
					if task_new.Type != types.BTN_TYPE_COMMAND {
						if task_new.Finished != 0 {
							driver.SetButtonLamp(task_new.Type, task_new.Floor, 0)
							button_lights[task_new.Floor][task_new.Type] = time.Since(time.Now())
						} else {
							driver.SetButtonLamp(task_new.Type, task_new.Floor, 1)
							button_lights[task_new.Floor][task_new.Type] = time.Now()
						}

					}

				} else {
					fmt.Println("AMANAGER: could not deserialize task from udp!")
				}

			default:
				fmt.Println("AMANAGER: message from udp unrecognisable!")
			}

		case <-time.After(types.PAUSE_AMAGER):
		}

		for i := 0; i < types.NUMBER_OF_FLOORS; i++ {
			for k := 0; k < 2; k++ {
				if time.Since(button_lights[i][k]) > types.TIMEOUT_LIGHT_ON {
					task_new.Type = k
					task_new.Floor = i
					task_new.Finished = 0
					task_new.Assigned = 255
					copy_assigned = nil
					for p := 0; p < len(assigned_tasks); p++ {
						copy_assigned = append(copy_assigned, assigned_tasks[p])
					}
					_, assigned_tasks = addTask(copy_assigned, task_new, elev_status)
				}
			}
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
