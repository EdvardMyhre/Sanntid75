package main

import (
	//"fmt"
	"time"
)

import (
	"./amanager"
	"./distributor"
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

	chan_elevStatus := make(chan types.Status) //Controller -> Amanager
	chan_elevTask := make(chan int)            //Amanager -> Controller
	//chan_networkToAmanager := make(chan types.MainData) //Network -> amanager
	chan_amanagerToNetwork := make(chan types.MainData) //Amanager -> network

	chan_distributorToNetwork := make(chan types.MainData) //Distributor -> network
	//chan_networkToDistributor := make(chan types.MainData) //network -> Distributor

	chan_backupReply := make(chan types.MainData) // all responses to backup -> amanager

	go pending_manager.Pending_task_manager(chan_button,
		chan_distStatus, chan_newDistOrder,
		chan_assignedTaskStatus, chan_assignedTask,
		chan_lostBackup)

	/*	go amanager.AssignedTasksManager(chan_elevStatus, chan_elevTask,
			chan_assignedTask, chan_assignedTaskStatus,
			chan_networkToAmanager, chan_amanagerToNetwork,
			chan_backupReply)

		go distributor.Task_distributor(chan_newDistOrder, chan_distStatus,
			chan_networkToDistributor, chan_distributorToNetwork)*/

	//AMANAGER OG DISTRIBUTOR UNDER ER FOR TEST VED DIREKTE KOBLING
	go amanager.AssignedTasksManager(chan_elevStatus, chan_elevTask,
		chan_assignedTask, chan_assignedTaskStatus,
		chan_distributorToNetwork, chan_amanagerToNetwork,
		chan_backupReply)

	go distributor.Task_distributor(chan_newDistOrder, chan_distStatus,
		chan_amanagerToNetwork, chan_distributorToNetwork)

	go elevator.Controller(chan_elevTask, chan_elevStatus)

	go elevator.ButtonPoller(chan_button)

	for {
		// select {
		// case <-chan_amanagerToNetwork:
		// case msg1 := <-chan_distributorToNetwork:
		// 	fmt.Println("	Distributor sent message to network:", msg1)

		// default:
		// }
		time.Sleep(time.Millisecond)
	}
}
