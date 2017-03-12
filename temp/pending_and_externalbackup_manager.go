
package main
import (
	"fmt"
	"time"
	"./Tester"
	"../Network"
	"../types"
)

func main(){
	
	//channel_distributor_to_network := make(chan Network.MainData)
	//channel_network_to_distributor := make(chan Network.MainData)
	
	//PENDING TASK MANAGER CHANNELS:
	//With button_intermediary
	channel_button_intermediary_to_pending_tasks_manager := make(chan Network.Button)
	//With distributor
	channel_distributor_to_pending_tasks_manager := make(chan Network.Button)
	channel_pending_tasks_manager_to_distributor := make(chan Network.Button)
	//With elevator handler
	channel_assigned_tasks_manager_to_pending_tasks_manager := make(chan Network.Button)
	channel_pending_tasks_manager_to_assigned_tasks_manager := make(chan Network.Button)
	//With Network Module
	channel_network_module_to_pending_tasks_manager := make(chan Network.MainData)
	channel_pending_tasks_manager_to_network_module := make(chan Network.MainData)
	
	go Pending_task_manager(	channel_button_intermediary_to_pending_tasks_manager,
								channel_distributor_to_pending_tasks_manager,			channel_pending_tasks_manager_to_distributor,
								channel_assigned_tasks_manager_to_pending_tasks_manager,		channel_pending_tasks_manager_to_assigned_tasks_manager,
								channel_network_module_to_pending_tasks_manager,		channel_pending_tasks_manager_to_network_module)

	go Tester.Test_pendingAndBackup_manager_buttonIntermediarySimulator(channel_button_intermediary_to_pending_tasks_manager)
	go Tester.Test_pendingandBackup_manager_assignedSimulator(channel_assigned_tasks_manager_to_pending_tasks_manager, channel_pending_tasks_manager_to_assigned_tasks_manager)
	fmt.Println("Main function for Pending tasks manager. Connectivity test")
	for{}

}


const(
	FLOORS = 4
)

type pendingListFloor struct {
	upOrder					uint8
	timestamp_upOrder		time.Time
	downOrder				uint8
	timestamp_downOrder 	time.Time
	internalOrder			uint8
	timestamp_internalOrder time.Time
	
}
	var pendingList [FLOORS+1] pendingListFloor

func Pending_task_manager(	channel_from_button_intermediary 	<-chan Network.Button,
							channel_from_distributor 			<-chan Network.Button, 					channel_to_distributor 				chan<- Network.Button,
							channel_from_assigned_tasks_manager <-chan Network.Button,					channel_to_assigned_tasks_manager 	chan<- Network.Button,
							channel_from_network 				<-chan Network.MainData,				channel_to_network 					chan<- Network.MainData) {
	fmt.Println("Pending Task Manager Go Routine startup")
	
	//Set up Pending Tasks Matrix

	//Zero out values just for added safety (Go should have done this automatically when initialized)
	for i := 0 ; i <=FLOORS ; i++{
		pendingList[i].upOrder = 0
		pendingList[i].downOrder = 0
		pendingList[i].internalOrder = 0
		//time.Time{} sets time value to 0 (aka reset value). Time library contains function to check if a time value is reset value.
		pendingList[i].timestamp_upOrder = time.Time{}
		pendingList[i].timestamp_downOrder = time.Time{}
		pendingList[i].timestamp_internalOrder = time.Time{}
	}
	
	fmt.Println("pendingList after being zeroed out: ", pendingList)
	
	for{
		//Behavior with button intermediary
		select{
			case message_buttonOrder := <- channel_from_button_intermediary:
				fmt.Println("Received new Order: ",message_buttonOrder)
				adjust_pendinglist(message_buttonOrder, false)
				
			default:
				//Do nothing
				
		}
			
		//Receive behavior with assigned tasks manager
		select{
			case message_buttonOrder := <- channel_from_assigned_tasks_manager:
				fmt.Println("Received assign Order: ",message_buttonOrder)
				adjust_pendinglist(message_buttonOrder, true)
					
			default:
					//Do nothing
		}
			
			
		//Receive behavior with task distributor
			//Controls variable which tells pending manager if it's supposed to give distributor new task
			
			
			
		//Receive behavior with backup-submodule
			//If receive anything, then merge lists (own routine for saving)
		
		
	}
	
	
	
	
	//Initiate External Backup Matrix
	//Initiate variables used for sending/receiving over channels
	
	//Implement behavior with button intermediary
	
	//Implement behavior with Distributor
	
	//Implement behavior with elevator handler
	
	//Implement behavior with network module

}


func adjust_pendinglist(adjust_pendinglist_message Network.Button, adjust_timestamp bool){
	if adjust_pendinglist_message.Floor >= len(pendingList) || adjust_pendinglist_message.Floor <= 0{
		//Illegal floor value
		fmt.Println("Pending adjustment: illegal FLOOR value")
		
	} else if adjust_pendinglist_message.Add != 0 {
		//Add order to pending
		if adjust_pendinglist_message.Button_type == types.BTN_TYPE_UP {
			fmt.Println("Pending adjustment: ADD UP ", adjust_pendinglist_message.Floor)
			pendingList[adjust_pendinglist_message.Floor].upOrder = 255
			if adjust_timestamp { pendingList[adjust_pendinglist_message.Floor].timestamp_upOrder = time.Now() }
		} else if adjust_pendinglist_message.Button_type == types.BTN_TYPE_DOWN {
			fmt.Println("Pending adjustment: ADD DOWN ", adjust_pendinglist_message.Floor)
			pendingList[adjust_pendinglist_message.Floor].downOrder = 255
			if adjust_timestamp { pendingList[adjust_pendinglist_message.Floor].timestamp_downOrder = time.Now() }
		} else if adjust_pendinglist_message.Button_type == types.BTN_TYPE_COMMAND {
			fmt.Println("Pending adjustment: ADD INTERNAL ", adjust_pendinglist_message.Floor)
			pendingList[adjust_pendinglist_message.Floor].internalOrder = 255
			if adjust_timestamp { pendingList[adjust_pendinglist_message.Floor].timestamp_internalOrder = time.Now() }
		}
							
							
	} else if adjust_pendinglist_message.Add == 0 {
		//Remove order from pending
		if adjust_pendinglist_message.Button_type == types.BTN_TYPE_UP {
			fmt.Println("Pending adjustment: REMOVE UP ", adjust_pendinglist_message.Floor)
			pendingList[adjust_pendinglist_message.Floor].upOrder = 0
			if adjust_timestamp { pendingList[adjust_pendinglist_message.Floor].timestamp_upOrder = time.Time{} }
		} else if adjust_pendinglist_message.Button_type == types.BTN_TYPE_DOWN {
			fmt.Println("Pending adjustment: REMOVE DOWN ", adjust_pendinglist_message.Floor)
			pendingList[adjust_pendinglist_message.Floor].downOrder = 0
			if adjust_timestamp { pendingList[adjust_pendinglist_message.Floor].timestamp_downOrder = time.Time{} }
		} else if adjust_pendinglist_message.Button_type == types.BTN_TYPE_COMMAND {
			fmt.Println("Pending adjustment: REMOVE INTERNAL ", adjust_pendinglist_message.Floor)
			pendingList[adjust_pendinglist_message.Floor].internalOrder = 0
			if adjust_timestamp { pendingList[adjust_pendinglist_message.Floor].timestamp_internalOrder = time.Time{} }
		}
	}
	
	//Prints below can be removed
	fmt.Println("Current pending list: ")
	for i := 1 ; i<len(pendingList);i++{
		fmt.Println(pendingList[i])
	}
	fmt.Println("")
}
