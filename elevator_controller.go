package main

import "./driver"
import "./types"
import "fmt"


func elevatorController(){

}


func main() {

	channel_elevator_controller_to_assigned_tasks_manager := make(chan types.Status)
	channel_assigned_tasks_manager_to_elevator_controller := make(chan types.Task)
	go elevatorController()
}
