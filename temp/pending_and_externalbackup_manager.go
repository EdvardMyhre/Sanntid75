
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
	gotPendingOrder 		uint8
	upOrder					uint8
	downOrder				uint8
	internalOrder			uint8
	timestamp_lastWorkedOn	time.Time
	
}


func Pending_task_manager(	channel_from_button_intermediary 	<-chan Network.Button,
							channel_from_distributor 			<-chan Network.Button, 					channel_to_distributor 				chan<- Network.Button,
							channel_from_assigned_tasks_manager <-chan Network.Button,					channel_to_assigned_tasks_manager 	chan<- Network.Button,
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
	
	//Initiate variables used for sending/receiving over button channels
	var message_buttonOrder Network.Button
	fmt.Println("So Go doesnt complain about variable usage: ",message_buttonOrder)
	
	
	for{
			//Behavior with button intermediary
			select{
				case message_buttonOrder := <- channel_from_button_intermediary:
					fmt.Println("Received new Order: ",message_buttonOrder)
					
					
					if message_buttonOrder.Floor >= len(pendingList) || message_buttonOrder.Floor <= 0{
						fmt.Println("Floor variable out of bounds")
					} else if message_buttonOrder.Button_type == types.BTN_TYPE_UP {
						fmt.Println("Command UP")
						pendingList[message_buttonOrder.Floor].upOrder = 255
					} else if message_buttonOrder.Button_type == types.BTN_TYPE_DOWN {
						fmt.Println("Command DOWN")
						pendingList[message_buttonOrder.Floor].downOrder = 255
					} else if message_buttonOrder.Button_type == types.BTN_TYPE_COMMAND {
						fmt.Println("Command INTERNAL")
						pendingList[message_buttonOrder.Floor].internalOrder = 255
					}
					fmt.Println("Curent pending list: ", pendingList)
					fmt.Println("")
					
				default:
					//Do nothing
				
			}
			
			//Receive behavior with assigned tasks manager
			select{
					case message_buttonOrder := <- channel_from_assigned_tasks_manager:
						fmt.Println("Received assign Order: ",message_buttonOrder)
					
						if message_buttonOrder.Floor >= len(pendingList) || message_buttonOrder.Floor <= 0{
							fmt.Println("Assigned: Floor variable out of bounds")
						} else if message_buttonOrder.Add != 0 {
							//Add order to pending
							if message_buttonOrder.Button_type == types.BTN_TYPE_UP {
								fmt.Println("Assigned Command UP")
								pendingList[message_buttonOrder.Floor].upOrder = 255
							} else if message_buttonOrder.Button_type == types.BTN_TYPE_DOWN {
								fmt.Println("Assigned Command DOWN")
								pendingList[message_buttonOrder.Floor].downOrder = 255
							} else if message_buttonOrder.Button_type == types.BTN_TYPE_COMMAND {
								fmt.Println("Assigned Command INTERNAL")
								pendingList[message_buttonOrder.Floor].internalOrder = 255
							}
							fmt.Println("Curent pending list: ", pendingList)
							fmt.Println("")							
							
						} else if message_buttonOrder.Add == 0 {
							//Remove order from pending
							if message_buttonOrder.Button_type == types.BTN_TYPE_UP {
								fmt.Println("Assigned Command REMOVE UP")
								pendingList[message_buttonOrder.Floor].upOrder = 0
							} else if message_buttonOrder.Button_type == types.BTN_TYPE_DOWN {
								fmt.Println("Assigned Command REMOVE DOWN")
								pendingList[message_buttonOrder.Floor].downOrder = 0
							} else if message_buttonOrder.Button_type == types.BTN_TYPE_COMMAND {
								fmt.Println("Assigned Command REMOVE INTERNAL")
								pendingList[message_buttonOrder.Floor].internalOrder = 0
							}
							fmt.Println("Curent pending list: ", pendingList)
							fmt.Println("")
						}
					default:
					//Do nothing
			}
			
			
			//Receive behavior with task distributor
			
			
			
			//Receive behavior with backup-submodule
		
		
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
