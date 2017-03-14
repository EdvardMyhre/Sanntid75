
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
	go Tester.Test_pendingandBackup_manager_networkSimulator(channel_network_module_to_backup_manager, channel_backup_manager_to_network_module)
	
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
	
	for{
//BEHAVIOR FOR RECEIVING FROM BUTTON POLLER
		select{
			case message_buttonOrder := <- channel_from_button_intermediary:
				adjust_pendinglist(message_buttonOrder.Type, message_buttonOrder.Floor, 0, true)			
			default:
				//Do nothing
				
		}
			
//BEHAVIOR FOR RECEIVING FROM ASSIGNED TASKS MANAGER
		select{
			case message_buttonOrder := <- channel_from_assigned_tasks_manager:
				adjust_pendinglist(message_buttonOrder.Type, message_buttonOrder.Floor, message_buttonOrder.Assigned, false)	
			default:
					//Do nothing
		}
			
			
//BEHAVIOR FOR RECEIVING FROM TASK DISTRIBUTOR
		select{
			case message_distributorStatus := <- channel_from_distributor:
				if message_distributorStatus.Finished == 0 {
					//Distributor is busy
					distributor_state.busy = 255
					distributor_state.timestamp_busyTime = time.Now()
				} else if message_distributorStatus.Finished != 0{
					//Distributor finished task, not busy anymore
					distributor_state.busy = 0
					distributor_state.timestamp_busyTime = time.Time{}
				}
				if message_distributorStatus.Assigned != 0 {
					adjust_pendinglist(message_distributorStatus.Type, message_distributorStatus.Floor, 0, true)
				}
			default:
				//Do nothing
		}
			
			
			
//BEHAVIOR FOR RECEIVING FROM BACKUP MANAGER
		select{
			case message_backup := <- channel_from_backup_manager:
				fmt.Println("Received message from backup manager. Must merge into pending list, and add timestamp to tasks", message_backup)
				for j := 0 ; j < len(message_backup.Data) ; j++ {
					if message_backup.Data[j][0] != types.BTN_TYPE_COMMAND {
						adjust_pendinglist(message_backup.Data[j][0], message_backup.Data[j][1], 255, true)
					}
				}	
			default:
		}
		
//BEHAVIOR FOR SENDING TO ASSIGNED TASK MANAGER		
		//Send behavior with assigned task manager. Loop through and look for either an internal order, or a timed out order.
		assigned_newOrderBool := false
		assigned_newOrderType := 0		
		for i := assigned_order_pendingList_startingIndex ; i <len(pendingList) ; i++ {
			//Check for internal order
			if (pendingList[i].internalOrder != 0) && (time.Time.IsZero(pendingList[i].timestamp_internalOrder) || (time.Since(pendingList[i].timestamp_internalOrder) > types.TIMEOUT_PENDINGLIST_ORDER)){
				
				assigned_newOrderBool = true
				assigned_newOrderType = types.BTN_TYPE_COMMAND
				assigned_order_pendingList_startingIndex = i
				break
			} else if (pendingList[i].upOrder != 0) && (time.Since(pendingList[i].timestamp_upOrder) > types.TIMEOUT_PENDINGLIST_ORDER) && !time.Time.IsZero(pendingList[i].timestamp_upOrder){
				assigned_newOrderBool = true
				assigned_newOrderType = types.BTN_TYPE_UP
				assigned_order_pendingList_startingIndex = i
				break
			} else if (pendingList[i].downOrder != 0) && (time.Since(pendingList[i].timestamp_downOrder) > types.TIMEOUT_PENDINGLIST_ORDER) && !time.Time.IsZero(pendingList[i].timestamp_downOrder) {
				assigned_newOrderBool = true
				assigned_newOrderType = types.BTN_TYPE_DOWN
				assigned_order_pendingList_startingIndex = i
			}
		}	
		//Do action if new order was found
		if assigned_newOrderBool {
			var sendOrder types.Task
			sendOrder.Floor = assigned_order_pendingList_startingIndex
			sendOrder.Type = assigned_newOrderType
			select{
				case channel_to_assigned_tasks_manager <- sendOrder:
					fmt.Println("SENT new order to assigned tasks manager: ",sendOrder)
					adjust_pendinglist(sendOrder.Type, sendOrder.Floor, 0, true)
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
					fmt.Println("SENT new order to assigned tasks manager FAILED DUE to TIMEOUT")
			}	
			if assigned_order_pendingList_startingIndex >= (len(pendingList)-1){
				assigned_order_pendingList_startingIndex = 0
			} else {
				assigned_order_pendingList_startingIndex += 1
			}	
		} else {
			assigned_order_pendingList_startingIndex = 0
		}
		

//BEHAVIOR FOR SENDING TO TASK DISTRIBUTOR
		distribute_newOrderBool := false
		distribute_newOrderType := 0
		distribute_newOrderIndex := 0
		
		if (distributor_state.busy == 0) || (time.Since(distributor_state.timestamp_busyTime) > types.TIMEOUT_MODULE_DISTRIBUTOR) {
			
			//Iterate through an look for unassigned UP and DOWN orders.
			//Send first one seen to distributor. Set busy to true	
			for i := 0 ; i <len(pendingList) ; i++ {
				if (pendingList[i].upOrder != 0) && (time.Time.IsZero(pendingList[i].timestamp_upOrder)) {
					distribute_newOrderBool = true
					distribute_newOrderType = types.BTN_TYPE_UP
					distribute_newOrderIndex = i
					break
				} else if (pendingList[i].downOrder != 0) && (time.Time.IsZero(pendingList[i].timestamp_downOrder)) {
					distribute_newOrderBool = true
					distribute_newOrderType = types.BTN_TYPE_DOWN
					distribute_newOrderIndex = i
					break
				}
			}	
		}
	
		if distribute_newOrderBool {
			var sendOrder types.Task
			sendOrder.Floor = distribute_newOrderIndex
			sendOrder.Type = distribute_newOrderType
			select {
				case channel_to_distributor <- sendOrder:
					fmt.Println("Message to Distributor SENT")
					adjust_pendinglist(sendOrder.Type, sendOrder.Floor, 0, true )
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
					fmt.Println("Message to Distributor FAILED to send")
			}		
		}
		time.Sleep(10*time.Millisecond)
	}
	
}




type struct_backup_element struct {
	BackupIP      string
	BackupData    [][]int
	
	
}

const (
	TEMPREQUEST = 1
	TEMPPUSH = 2
)



func Backup_manager(	channel_from_network 	<-chan types.MainData,		channel_to_network 			chan<- types.MainData,
																			channel_to_pending_manager 	chan<- types.MainData){
	//var message_maindata types.MainData
//CREATE MATRIX
	var backup_matrix []struct_backup_element
	var sendQueue_push []string
	var sendQueue_request []string
	var sendQueue_pending_manager []string

	for {
//BEHAVIOR FOR RECEIVING FROM NETWORK
		select {
			case network_message := <- channel_from_network:
				fmt.Println("Received something from network: ",network_message)
		//"REQUEST BACKUP" RECEIVED
				if network_message.Type == TEMPREQUEST {
					var request_command_already_exists bool
					request_command_already_exists = false
					for i := 0 ; i < len(sendQueue_request) ; i++ {
						if sendQueue_request[i] == network_message.Source {
							request_command_already_exists = true
							break
						}
					}
					if request_command_already_exists == false {
						sendQueue_request = append(sendQueue_request, network_message.Source)
						fmt.Println("Received NEW REQUEST command message. Added to queue. Current request queue: ", sendQueue_request)
					}
				

		//"PUSH BACKUP" RECEIVED		
				} else if network_message.Type == TEMPPUSH {
					var backup_already_exists bool
					var push_command_already_exists bool
					backup_already_exists = false
					push_command_already_exists = false
					//Check if a backup already exists
					for i := 0 ; i < len(backup_matrix) ; i++ {
						if backup_matrix[i].BackupIP == network_message.Source {
							backup_matrix[i].BackupData = network_message.Data
							backup_already_exists = true
							break
						}
					}
					//if not, create new element and append to matrix
					if backup_already_exists == false {
						var new_backup struct_backup_element
						new_backup.BackupIP = network_message.Source
						new_backup.BackupData = network_message.Data
						backup_matrix = append(backup_matrix, new_backup)
					}
					//Check if response command already exists
					for i:= 0 ; i < len(sendQueue_push) ; i++ {
						if sendQueue_push[i] == network_message.Source {
							push_command_already_exists = true
						}
					}
					
					//if not, create new command and append to list
					if push_command_already_exists == false {
						sendQueue_push = append(sendQueue_push, network_message.Source)
						fmt.Println("Received NEW push command message. Added to queue. Current queue: ", sendQueue_push)
					}
					
		//"BACKUP LOST" RECEIVED	
				} else if network_message.Type == types.BACKUP_LOST {
					//Check if task is already in queue
					var command_already_exists bool
					command_already_exists = false
					for i := 0 ; i< len(sendQueue_pending_manager) ; i++ {
						if sendQueue_pending_manager[i] == network_message.Source {
							command_already_exists = true
							fmt.Println("Received an already exsting backup lost message. Duplicate NOT added to queue. Current queue: ",sendQueue_pending_manager)
						}
					}
					if command_already_exists == false {
						sendQueue_pending_manager = append(sendQueue_pending_manager, network_message.Source)
						fmt.Println("Received new backup lost message. Current queue: ",sendQueue_pending_manager)
					}
				}
			
			default:
			//Do nothing
		}
		
//BEHAVIOR FOR SENDING "PUSH" RESPONSE TO NETWORK
		//If there's a an order in network sendQueue
		if len(sendQueue_push) > 0 {
			var push_index int
			var push_backup_exists bool
			push_backup_exists = false
			for i := 0 ; i<len(backup_matrix) ; i++ {
				if backup_matrix[i].BackupIP == sendQueue_push[0] {
					push_index = i
					push_backup_exists = true
					break
				}
			}
			
			var push_message types.MainData
			push_message.Destination = sendQueue_push[0]
			if push_backup_exists {
				push_message.Data = backup_matrix[push_index].BackupData
			}
			select{
				case channel_to_network <- push_message:
					fmt.Println("Sent push response back to: ",sendQueue_push[0])
					//Delete element in index 0
					sendQueue_push = append(sendQueue_push[:0])
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
			}
			
		}
//BEHAVIOR FoR SENDING "REQUEST" RESPONSE TO NETWORK		
		if len(sendQueue_request) > 0 {
			
		}
		
//BEHAVIor FOR SENDING "BACKUP LOST" TO PENDING MANAGER
		if len(sendQueue_pending_manager) > 0 {
			
		}
	}
																				
}



//FLOOR, ASSIGNED, BUTTON TYPE, TIMESTAMP
func adjust_pendinglist(adjust_type int, adjust_floor int,adjust_assigned int, adjust_timestamp bool){
	if adjust_floor >= len(pendingList) || adjust_floor < 0{
		//Illegal floor value
		fmt.Println("Pending adjustment: illegal FLOOR value")
		
	} else if adjust_assigned == 0 {
		//Un-assigned order seen, add to pending list. Set timestamp if adjust_timestamp is true
		if adjust_type == types.BTN_TYPE_UP {
			fmt.Println("Pending adjustment: ADD UP ", adjust_floor)
			pendingList[adjust_floor].upOrder = 255
			if adjust_timestamp { pendingList[adjust_floor].timestamp_upOrder = time.Now() }
		} else if adjust_type == types.BTN_TYPE_DOWN {
			fmt.Println("Pending adjustment: ADD DOWN ", adjust_floor)
			pendingList[adjust_floor].downOrder = 255
			if adjust_timestamp { pendingList[adjust_floor].timestamp_downOrder = time.Now() }
		} else if adjust_type == types.BTN_TYPE_COMMAND {
			fmt.Println("Pending adjustment: ADD INTERNAL ", adjust_floor)
			pendingList[adjust_floor].internalOrder = 255
			if adjust_timestamp { pendingList[adjust_floor].timestamp_internalOrder = time.Now() }
		}
							
							
	} else if adjust_assigned != 0 {
		//Order has been assigned. No longer pending. Remove from pending list. Remove timestamp if adjust_timestamp is true
		if adjust_type == types.BTN_TYPE_UP {
			fmt.Println("Pending adjustment: REMOVE UP ", adjust_floor)
			pendingList[adjust_floor].upOrder = 0
			if adjust_timestamp { pendingList[adjust_floor].timestamp_upOrder = time.Time{} }
		} else if adjust_type == types.BTN_TYPE_DOWN {
			fmt.Println("Pending adjustment: REMOVE DOWN ", adjust_floor)
			pendingList[adjust_floor].downOrder = 0
			if adjust_timestamp { pendingList[adjust_floor].timestamp_downOrder = time.Time{} }
		} else if adjust_type == types.BTN_TYPE_COMMAND {
			fmt.Println("Pending adjustment: REMOVE INTERNAL ", adjust_floor)
			pendingList[adjust_floor].internalOrder = 0
			if adjust_timestamp { pendingList[adjust_floor].timestamp_internalOrder = time.Time{} }
		}
	}
	
	//Prints below can be removed
	//fmt.Println("Current pending list: ")
	//for i := 0 ; i<len(pendingList);i++{
	//	fmt.Println(pendingList[i])
	//}
	//fmt.Println("")
}
