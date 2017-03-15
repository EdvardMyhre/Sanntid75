package main

import (
	//"fmt"
	"time"
)

import (
	"./amanager"
	"./backup_manager"
	"./distributor"
	"./driver"
	"./elevator"
	"./network_module"
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

	chan_distributorToNetwork := make(chan types.MainData) //Distributor -> network
	chan_networkToDistributor := make(chan types.MainData) //network -> Distributor

	chan_backupReply := make(chan types.MainData) // Backup -> Amanager

	chan_backupToNetwork := make(chan types.MainData) //Backup -> Network
	chan_networkToBackup := make(chan types.MainData) //Network -> Backup

	go network.Network_start(chan_networkToDistributor, chan_networkToBackup, chan_networkToAmanager,
		chan_distributorToNetwork, chan_backupToNetwork, chan_amanagerToNetwork, chan_backupReply)

	time.Sleep(time.Second * 5)

	go pending_manager.Pending_task_manager(chan_button,
		chan_distStatus, chan_newDistOrder,
		chan_assignedTaskStatus, chan_assignedTask,
		chan_lostBackup)

	go backup_manager.Backup_manager(chan_networkToBackup, chan_backupToNetwork, chan_lostBackup)

	go amanager.Assigned_tasks_manager(chan_elevStatus, chan_elevTask,
		chan_assignedTask, chan_assignedTaskStatus,
		chan_networkToAmanager, chan_amanagerToNetwork,
		chan_backupReply)

	go distributor.Task_distributor(chan_newDistOrder, chan_distStatus,
		chan_networkToDistributor, chan_distributorToNetwork)

	go elevator.Controller(chan_elevTask, chan_elevStatus)

	go elevator.Button_poller(chan_button)

	for {

		time.Sleep(time.Millisecond)
	}
}
