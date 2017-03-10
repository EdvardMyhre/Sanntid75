
package main
import (
	"fmt"
	"time"
	"./Tester"
	"../Network"
)

func main(){
	
	//channel_distributor_to_network := make(chan Network.MainData)
	//channel_network_to_distributor := make(chan Network.MainData)
	
	//PENDING TASK MANAGER CHANNELS:
	//With button_intermediary
	channel_button_intermediary_to_pending_tasks_manager := make(chan Network.Button)
	//With distributor
	channel_distributor_to_pending_tasks_manager := make(chan Network.InternalMessage)
	channel_pending_tasks_manager_to_distributor := make(chan Network.InternalMessage)
	//With elevator handler
	channel_assigned_tasks_manager_to_pending_tasks_manager := make(chan Network.InternalMessage)
	channel_pending_tasks_manager_to_assigned_tasks_manager := make(chan Network.InternalMessage)
	//With Network Module
	channel_network_module_to_pending_tasks_manager := make(chan Network.MainData)
	channel_pending_tasks_manager_to_network_module := make(chan Network.MainData)
	
	go Pending_task_manager(	channel_button_intermediary_to_pending_tasks_manager,
								channel_distributor_to_pending_tasks_manager,			channel_pending_tasks_manager_to_distributor,
								channel_assigned_tasks_manager_to_pending_tasks_manager,		channel_pending_tasks_manager_to_assigned_tasks_manager,
								channel_network_module_to_pending_tasks_manager,		channel_pending_tasks_manager_to_network_module)

	go Tester.Test_pendingAndBackup_manager_buttonIntermediarySimulator(channel_button_intermediary_to_pending_tasks_manager)							
	fmt.Println("Main function for Pending tasks manager. Connectivity test")
	for{}

}


const(
	FLOORS = 4
)

type pendingListFloor struct {
	gotPendingOrder 		uint8
	upOrder					uint8
	downOrder				uint8
	internalOrder			uint8
	timestamp_lastWorkedOn	time.Time
	
}


func Pending_task_manager(	channel_from_button_intermediary 	<-chan Network.Button,
							channel_from_distributor 			<-chan Network.InternalMessage, 		channel_to_distributor 				chan<- Network.InternalMessage,
							channel_from_assigned_tasks_manager <-chan Network.InternalMessage,			channel_to_assigned_tasks_manager 	chan<- Network.InternalMessage,
							channel_from_network 				<-chan Network.MainData,				channel_to_network 					chan<- Network.MainData) {
	fmt.Println("Pending Task Manager Go Routine startup")
	
	//Set up Pending Tasks Matrix
	var pendingList [FLOORS+1] pendingListFloor
	//Zero out values jsut for added safety (Go should have done this automatically when initialized)
	for i := 0 ; i <=FLOORS ; i++{
		pendingList[i].gotPendingOrder = 0
		pendingList[i].upOrder = 0
		pendingList[i].downOrder = 0
		pendingList[i].internalOrder = 0
		pendingList[i].timestamp_lastWorkedOn = time.Time{}
	}
	
	fmt.Println("pendingList after being zeroed out: ", pendingList)
	
	//Initiate variables used for sending/receiving over channels
	var message_new_buttonOrder Network.Button
	fmt.Println("So Go doesnt complain about variable usage: ",message_new_buttonOrder)
	
	
	for{
			//Behavior with button intermediary
			//case message_distributingOrder := <- channel_from_network:
			select{
				case message_new_buttonOrder := <- channel_from_button_intermediary:
					fmt.Println("Received something: ",message_new_buttonOrder)
				default:
					//Do nothing
				
			}
		
		
	}
	
	
	
	
	//Initiate External Backup Matrix
	//Initiate variables used for sending/receiving over channels
	
	//Implement behavior with button intermediary
	
	//Implement behavior with Distributor
	
	//Implement behavior with elevator handler
	
	//Implement behavior with network module
	
	
	for{
		//
	}

}
