package main

import (
	//"fmt"
	"time"
)

import (
	"./amanager"
	"./driver"
	"./elevator"
	"./pending_manager"
	"./types"
)

func main() {
	driver.Init()
	chan_button := make(chan types.Task) //Poller -> Pending

	chan_distStatus := make(chan types.Task)   //Dist -> Pending
	chan_newDistOrder := make(chan types.Task) //Pending -> Dist

	chan_assignedTaskStatus := make(chan types.Task) //Assigned -> Pending
	chan_assignedTask := make(chan types.Task)       //Pending -> Assigned

	chan_lostBackup := make(chan types.MainData) //Backup -> Pending

	chan_elevStatus := make(chan types.Status)          //Controller -> Amanager
	chan_elevTask := make(chan int)                     //Amanager -> Controller
	chan_networkToAmanager := make(chan types.MainData) //Network -> amanager
	chan_amanagerToNetwork := make(chan types.MainData) //Amanager -> network

	go pending_manager.Pending_task_manager(chan_button,
		chan_distStatus, chan_newDistOrder,
		chan_assignedTaskStatus, chan_assignedTask,
		chan_lostBackup)

	go amanager.AssignedTasksManager(chan_elevStatus, chan_elevTask,
		chan_assignedTask, chan_assignedTaskStatus,
		chan_networkToAmanager, chan_amanagerToNetwork)

	go elevator.Controller(chan_elevTask, chan_elevStatus)

	go elevator.ButtonPoller(chan_button)

	for {
		select {
		case <-chan_amanagerToNetwork:
		default:
		}
		time.Sleep(time.Millisecond)
	}
}
