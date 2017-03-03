
package main
import (
	"fmt"
	"time"
	"./Tester"
	"../Network"
)

func main(){
	
	channel_distributor_to_network := make(chan Network.MainData)
	channel_network_to_distributor := make(chan Network.MainData)
	
	channel_task_manager_to_distributor := make(chan Network.NewOrder)
	channel_distributor_to_task_manager := make(chan Network.NewOrder)
	
	go task_distributor(channel_distributor_to_task_manager, channel_task_manager_to_distributor,
						channel_distributor_to_network, channel_network_to_distributor)
						
	//forever loop so the main function does not terminate
	go Tester.Test_distributor_networkEmulator(channel_network_to_distributor, channel_distributor_to_network)
	
	go Tester.Test_distributor_taskmanagerEmulator(channel_task_manager_to_distributor, channel_distributor_to_task_manager)
	
	for{}

}


func task_distributor(channel_to_task_manager chan Network.NewOrder, channel_from_task_manager chan Network.NewOrder,
						channel_to_network chan Network.MainData, channel_from_network chan Network.MainData) {
	fmt.Println("task distributor, Connectivity test")
	const (
		waiting_for_task = 1		//001
		waiting_for_response = 2	//010
		send_choice = 5				//100
									//111
	)
	//State machine variable
	var task_distributor_state int = waiting_for_task
	
	//Timestamp variable
	task_dist_timestamp := time.Now()
	
	//Variables used for channels
	var message_distributor_to_network Network.MainData
	var message_distributor_from_network Network.MainData
	var message_distributor_to_task_manager Network.NewOrder
	var message_distributor_from_task_manager Network.NewOrder
	
	fmt.Println("Just so Go doesnt complain about varible usage: ",message_distributor_to_network, message_distributor_from_network, message_distributor_to_task_manager, message_distributor_from_task_manager )
	
	//TEMP: Writing out starting timestamp
	fmt.Println("Timestamp: ",task_dist_timestamp)
	
	
	
	
	
	for {
		switch task_distributor_state{
		case waiting_for_task:
			//Waits for new task from task_manager
			//fmt.Println("Waiting for task")
			
			select{
				case message_distributor_from_task_manager := <-channel_from_task_manager:
					if message_distributor_from_task_manager.Message_type == Network.ID_MSG_TYPE_DISTIBUTOR_NEW_COMMAND {
						fmt.Println("Distributor received new command from Task Manager")
						fmt.Println("Send request message to network, and enter waiting_for_response state")
					
						task_distributor_state = waiting_for_response
					} else {
						fmt.Println("Received some other message")
					}
					
				default:
					//Do nothing for now.

			
			}
			
			
		case waiting_for_response:
			//Waits for other nodes to respond
			//fmt.Println("Waiting for response")
			task_distributor_state = 0
			
			
		case send_choice:
			//Sends message containing the node chosen for task
			
			
			
		default:
			task_distributor_state = waiting_for_task	
			
			fmt.Println("Time difference: ",time.Since(task_dist_timestamp).Seconds())
			
		}
	}
	
}
