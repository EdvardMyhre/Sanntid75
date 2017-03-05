
package main
import (
	"fmt"
	//"time"
	//"./Tester"
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
	channel_elevator_handler_to_pending_tasks_manager := make(chan Network.InternalMessage)
	channel_pending_tasks_manager_to_elevator_handler := make(chan Network.InternalMessage)
	//With Network Module
	channel_network_module_to_pending_tasks_manager := make(chan Network.MainData)
	channel_pending_tasks_manager_to_network_module := make(chan Network.MainData)
	
	go Pending_task_manager(	channel_button_intermediary_to_pending_tasks_manager,
								channel_distributor_to_pending_tasks_manager,			channel_pending_tasks_manager_to_distributor,
								channel_elevator_handler_to_pending_tasks_manager,		channel_pending_tasks_manager_to_elevator_handler,
								channel_network_module_to_pending_tasks_manager,		channel_pending_tasks_manager_to_network_module)
	
	//go task_distributor(channel_distributor_to_task_manager, channel_task_manager_to_distributor,
	//					channel_distributor_to_network, channel_network_to_distributor)
						
	
	//go Tester.Test_distributor_networkEmulator(channel_network_to_distributor, channel_distributor_to_network)
	
	//go Tester.Test_distributor_taskmanagerEmulator(channel_task_manager_to_distributor, channel_distributor_to_task_manager)
	
	//forever loop so the main function does not terminate
	fmt.Println("Main function for Pending tasks manager. Connectivity test")
	for{}

}





func Pending_task_manager(	channel_from_button_intermediary 	<-chan Network.Button,
							channel_from_distributor 			<-chan Network.InternalMessage, 		channel_to_distributor 		chan<- Network.InternalMessage,
							channel_from_elevator_handler 		<-chan Network.InternalMessage,			channel_to_elevator_handler chan<- Network.InternalMessage,
							channel_from_network 				<-chan Network.MainData,				channel_to_network 			chan<- Network.MainData) {
	fmt.Println("Pending Task Manager Go Routine startup")
	//Set up Pending Tasks Matrix
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
