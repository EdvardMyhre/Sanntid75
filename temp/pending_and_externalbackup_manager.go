
package main
import (
	"fmt"
	"time"
	"./Tester"
	"../types"
)

func main(){
	
	//channel_distributor_to_network := make(chan Network.MainData)
	//channel_network_to_distributor := make(chan Network.MainData)
	
	//PENDING TASK MANAGER CHANNELS:
	//With button_intermediary
	channel_button_intermediary_to_pending_tasks_manager := make(chan types.Button)
	//With distributor
	channel_distributor_to_pending_tasks_manager := make(chan types.Task)
	channel_pending_tasks_manager_to_distributor := make(chan types.Task)
	//With assigned tasks manager
	channel_assigned_tasks_manager_to_pending_tasks_manager := make(chan types.Task)
	channel_pending_tasks_manager_to_assigned_tasks_manager := make(chan types.Task)
	//With backup manager
	channel_backup_manager_to_pending_manager := make(chan types.MainData)
	
	//Backup manager communication with network module
	channel_network_module_to_backup_manager := make(chan types.MainData)
	channel_backup_manager_to_network_module := make(chan types.MainData)
	
	
	//Module go routines
	go Pending_task_manager(	channel_button_intermediary_to_pending_tasks_manager,
								channel_distributor_to_pending_tasks_manager,			channel_pending_tasks_manager_to_distributor,
								channel_assigned_tasks_manager_to_pending_tasks_manager,		channel_pending_tasks_manager_to_assigned_tasks_manager,
								channel_backup_manager_to_pending_manager)
	
	go Backup_manager(			channel_network_module_to_backup_manager,	channel_backup_manager_to_network_module,
																			channel_backup_manager_to_pending_manager)
	
	//Test functions, Can be removed when combining
	go Tester.Test_pendingAndBackup_manager_buttonIntermediarySimulator(channel_button_intermediary_to_pending_tasks_manager)
	go Tester.Test_pendingandBackup_manager_assignedSimulator(channel_assigned_tasks_manager_to_pending_tasks_manager, channel_pending_tasks_manager_to_assigned_tasks_manager)
	go Tester.Test_pendingandBackup_manager_distributorSimulator(channel_distributor_to_pending_tasks_manager, channel_pending_tasks_manager_to_distributor)
	
	
	fmt.Println("Main function for Pending tasks manager. Connectivity test")
	for{}

}



type struct_pendingListFloor struct {
	upOrder					uint8
	timestamp_upOrder		time.Time
	downOrder				uint8
	timestamp_downOrder 	time.Time
	internalOrder			uint8
	timestamp_internalOrder time.Time
	
}
	var pendingList [types.NUMBER_OF_FLOORS] struct_pendingListFloor
	
	
type struct_distributor_state struct {
	busy					uint8
	timestamp_busyTime		time.Time
}
	var distributor_state struct_distributor_state

	
func Pending_task_manager(	channel_from_button_intermediary 	<-chan types.Button,
							channel_from_distributor 			<-chan types.Task, 						channel_to_distributor 				chan<- types.Task,
							channel_from_assigned_tasks_manager <-chan types.Task	,					channel_to_assigned_tasks_manager 	chan<- types.Task,
							channel_from_backup_manager			<-chan types.MainData,																			) {
	fmt.Println("Pending Task Manager Go Routine startup")
	
	//Set up Pending Tasks Matrix

	//Zero out values just for added safety (Go should have done this automatically when initialized)
	for i := 0 ; i <len(pendingList) ; i++{
		pendingList[i].upOrder = 0
		pendingList[i].downOrder = 0
		pendingList[i].internalOrder = 0
		//time.Time{} sets time value to 0 (aka reset value). Time library contains function to check if a time value is reset value.
		pendingList[i].timestamp_upOrder = time.Time{}
		pendingList[i].timestamp_downOrder = time.Time{}
		pendingList[i].timestamp_internalOrder = time.Time{}
	}
	distributor_state.busy = 0
	distributor_state.timestamp_busyTime = time.Time{}
	
	//Variable declarations
	assigned_order_pendingList_startingIndex := 0
	
	fmt.Println("pendingList after being zeroed out: ", pendingList)
	
	for{
		//Receive Behavior with button intermediary
		select{
			case message_buttonOrder := <- channel_from_button_intermediary:
				//fmt.Println("Received new Order: ",message_buttonOrder)
				var buttontoTaskConvert types.Task
				buttontoTaskConvert.Floor = message_buttonOrder.Floor
				buttontoTaskConvert.Type = message_buttonOrder.Type
				buttontoTaskConvert.Assigned = 0
				adjust_pendinglist(buttontoTaskConvert, false)
				
			default:
				//Do nothing
				
		}
			
		//Receive Behavior with assigned tasks manager
		select{
			case message_buttonOrder := <- channel_from_assigned_tasks_manager:
				//fmt.Println("Received assign Order: ",message_buttonOrder)
				adjust_pendinglist(message_buttonOrder, true)
					
			default:
					//Do nothing
		}
			
			
		//Receive behavior with task distributor
			//Controls variable which tells pending manager if it's supposed to give distributor new task or not (if distributor is busy or not)
		select{
			case message_distributorStatus := <- channel_from_distributor:
				if message_distributorStatus.Finished == 0 {
					//Distributor is busy
					distributor_state.busy = 255
					distributor_state.timestamp_busyTime = time.Now()
					//fmt.Println("Received distributor state BUSY :" , distributor_state, "Finished value:", message_distributorStatus.Finished)
				} else if message_distributorStatus.Finished != 0{
					//Distributor finished task, not 
					distributor_state.busy = 0
					distributor_state.timestamp_busyTime = time.Time{}
					//fmt.Println("Received distributor state READY :" , distributor_state, "Finished value:", message_distributorStatus.Finished)
				}
			default:
				//Do nothing
		}
			
			
			
		//Receive behavior with backup-submodule
			//If receive anything, then merge lists (own routine for saving)
		select{
			case message_backup := <- channel_from_backup_manager:
				fmt.Println("Received message from backup manager. Must merge into pending list, and add timestamp to tasks", message_backup)
				fmt.Println("")
				fmt.Println("")
				fmt.Println("")
				//Pending only receives something here if the computer has lost connection with an elevator who used it as backup
				//Loop through it and add all tasks in the backup matrix to pending matrix, with timestamp.
			default:
		}
		
		
		//Send behavior with assigned task manager. Loop through and look for either an internal order, or a timed out order.
		assigned_newOrderBool := false
		assigned_newOrderType := 0
		//assigned_order_pendingList_startingIndex := 1		//This variable is set outside for loop
		
		for i := assigned_order_pendingList_startingIndex ; i <len(pendingList) ; i++ {
			//Check for internal order
			if (pendingList[i].internalOrder != 0) && (time.Time.IsZero(pendingList[i].timestamp_internalOrder) || (time.Since(pendingList[i].timestamp_internalOrder) > types.TIMEOUT_PENDINGLIST_ORDER)){
				
				assigned_newOrderBool = true
				assigned_newOrderType = types.BTN_TYPE_COMMAND
				assigned_order_pendingList_startingIndex = i
				break
			} else if (pendingList[i].upOrder != 0) && (time.Since(pendingList[i].timestamp_upOrder) > types.TIMEOUT_PENDINGLIST_ORDER){
				assigned_newOrderBool = true
				assigned_newOrderType = types.BTN_TYPE_UP
				assigned_order_pendingList_startingIndex = i
				break
			} else if (pendingList[i].downOrder != 0) && (time.Since(pendingList[i].timestamp_downOrder) > types.TIMEOUT_PENDINGLIST_ORDER) {
				assigned_newOrderBool = true
				assigned_newOrderType = types.BTN_TYPE_DOWN
				assigned_order_pendingList_startingIndex = i
			}
			//Check for up order
			
			//Check for down order
		}
		
		//Do action if new order was found
		if assigned_newOrderBool {
			var sendOrder types.Task
			sendOrder.Floor = assigned_order_pendingList_startingIndex
			sendOrder.Type = assigned_newOrderType
			select{
				case channel_to_assigned_tasks_manager <- sendOrder:
					fmt.Println("                                                    SENT new order to assigned tasks manager: ",sendOrder)
					if assigned_newOrderType == types.BTN_TYPE_COMMAND {
						pendingList[assigned_order_pendingList_startingIndex].timestamp_internalOrder = time.Now()
					} else if assigned_newOrderType == types.BTN_TYPE_UP {
						pendingList[assigned_order_pendingList_startingIndex].timestamp_upOrder = time.Now()
					} else if assigned_newOrderType == types.BTN_TYPE_DOWN {
						pendingList[assigned_order_pendingList_startingIndex].timestamp_downOrder = time.Now()
					}	
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
					fmt.Println("                                                    SENT new order to assigned tasks manager FAILED DUE to TIMEOUT")
			}
			
			
			
			
			
			if assigned_order_pendingList_startingIndex >= (len(pendingList)-1){
				assigned_order_pendingList_startingIndex = 0
			} else {
				assigned_order_pendingList_startingIndex += 1
			}
			
			
		} else {
			assigned_order_pendingList_startingIndex = 0
		}
		
		
	}
	
	
	
	
	//Initiate External Backup Matrix
	//Initiate variables used for sending/receiving over channels
	
	//Implement behavior with button intermediary
	
	//Implement behavior with Distributor
	
	//Implement behavior with elevator handler
	
	//Implement behavior with network module

}


func Backup_manager(	channel_from_network 	<-chan types.MainData,		channel_to_network 			chan<- types.MainData,
																			channel_to_pending_manager 	chan<- types.MainData){
	var message_maindata types.MainData
	for{
		time.Sleep(15*time.Second)
		channel_to_pending_manager <- message_maindata
	}
																				
}
func adjust_pendinglist(adjust_pendinglist_message types.Task, adjust_timestamp bool){
	if adjust_pendinglist_message.Floor >= len(pendingList) || adjust_pendinglist_message.Floor < 0{
		//Illegal floor value
		fmt.Println("Pending adjustment: illegal FLOOR value")
		
	} else if adjust_pendinglist_message.Assigned == 0 {
		//Un-assigned order seen, add to pending list. Set timestamp if adjust_timestamp is true
		if adjust_pendinglist_message.Type == types.BTN_TYPE_UP {
			fmt.Println("Pending adjustment: ADD UP ", adjust_pendinglist_message.Floor)
			pendingList[adjust_pendinglist_message.Floor].upOrder = 255
			if adjust_timestamp { pendingList[adjust_pendinglist_message.Floor].timestamp_upOrder = time.Now() }
		} else if adjust_pendinglist_message.Type == types.BTN_TYPE_DOWN {
			fmt.Println("Pending adjustment: ADD DOWN ", adjust_pendinglist_message.Floor)
			pendingList[adjust_pendinglist_message.Floor].downOrder = 255
			if adjust_timestamp { pendingList[adjust_pendinglist_message.Floor].timestamp_downOrder = time.Now() }
		} else if adjust_pendinglist_message.Type == types.BTN_TYPE_COMMAND {
			fmt.Println("Pending adjustment: ADD INTERNAL ", adjust_pendinglist_message.Floor)
			pendingList[adjust_pendinglist_message.Floor].internalOrder = 255
			if adjust_timestamp { pendingList[adjust_pendinglist_message.Floor].timestamp_internalOrder = time.Now() }
		}
							
							
	} else if adjust_pendinglist_message.Assigned != 0 {
		//Order has been assigned. No longer pending. Remove from pending list. Remove timestamp if adjust_timestamp is true
		if adjust_pendinglist_message.Type == types.BTN_TYPE_UP {
			fmt.Println("Pending adjustment: REMOVE UP ", adjust_pendinglist_message.Floor)
			pendingList[adjust_pendinglist_message.Floor].upOrder = 0
			if adjust_timestamp { pendingList[adjust_pendinglist_message.Floor].timestamp_upOrder = time.Time{} }
		} else if adjust_pendinglist_message.Type == types.BTN_TYPE_DOWN {
			fmt.Println("Pending adjustment: REMOVE DOWN ", adjust_pendinglist_message.Floor)
			pendingList[adjust_pendinglist_message.Floor].downOrder = 0
			if adjust_timestamp { pendingList[adjust_pendinglist_message.Floor].timestamp_downOrder = time.Time{} }
		} else if adjust_pendinglist_message.Type == types.BTN_TYPE_COMMAND {
			fmt.Println("Pending adjustment: REMOVE INTERNAL ", adjust_pendinglist_message.Floor)
			pendingList[adjust_pendinglist_message.Floor].internalOrder = 0
			if adjust_timestamp { pendingList[adjust_pendinglist_message.Floor].timestamp_internalOrder = time.Time{} }
		}
	}
	
	//Prints below can be removed
	fmt.Println("Current pending list: ")
	for i := 0 ; i<len(pendingList);i++{
		fmt.Println(pendingList[i])
	}
	fmt.Println("")
}
