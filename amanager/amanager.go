package amanager

import "./driver"
import "./types"
import "./elevator"

import "fmt"
import "time"

func AssignedTasksManager(elev_status_c <-chan types.Status, elev_task_c chan<- int, statusc <-chan types.Task, taskc chan<- types.Task, udp_rx_c <-chan types.MainData, udp_tx_c chan<- types.MainData){


	//Initializing variables
	var assigned_tasks []types.Task
	var msg types.MainData{}


	//Boot routine
	//Fetch backup if available, network module has already found our backup or created a new
	msg.Destination = "backup"
	msg.Message_type = MESSAGE_TYPE_REQUEST_BACKUP
	udp_tx_c <- msg
	select{
		case msg := <- udp_rx_c:
			assigned_tasks = slice2tasks(msg.Data)
		case <- time.After(time.Second * 10):
	}



}





//Functions

func slice2tasks(slice [][]int) []types.Task{
	l := len(slice)
	tasks := make([]types.Task, l)
	for i := 0; i < len(slice); i++{
		tasks[i].Type = slice[i][0]
		tasks[i].Floor = slice[i][1]
		tasks[i].Add = slice[i][2] 
	}
	return tasks
}

func tasks2slice(tasks []types.Task) [][]int{
	l = len(tasks)
	slice := make([][]int, l, 3)
	for i := 0; i < len(tasks); i++{
		slice[i][0] = tasks[i].Type
		slice[i][1] = tasks[i].Floor
		slice[i][2] = tasks[i].Add
	}
	return slice
}

func addTask(tasks []types.Task, task types.Task, status types.Status) (weight int, queuePos int){
	//task ligger i tasks fra før av kan jeg anta. Skal nå få en vekt på å utføre task.
	weight := 0
	direction := UP
	if status.Prev_floor > status.Floor || status.Floor = NUMBER_OF_FLOORS-1{
		direction = DOWN
	}


}

const(
	UP = 1
	DOWN = 0
)