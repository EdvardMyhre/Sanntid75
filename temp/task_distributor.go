
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
	
	channel_task_manager_to_distributor := make(chan Network.InternalMessage)
	channel_distributor_to_task_manager := make(chan Network.InternalMessage)
	
	go task_distributor(channel_distributor_to_task_manager, channel_task_manager_to_distributor,
						channel_distributor_to_network, channel_network_to_distributor)
						
	
	go Tester.Test_distributor_networkEmulator(channel_network_to_distributor, channel_distributor_to_network)
	
	go Tester.Test_distributor_taskmanagerEmulator(channel_task_manager_to_distributor, channel_distributor_to_task_manager)
	
	//forever loop so the main function does not terminate
	for{}

}


func task_distributor(channel_to_task_manager chan Network.InternalMessage, channel_from_task_manager chan Network.InternalMessage,
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
	
	//currentorder = {floor, direction}
	var currentOrder []int
	//Response array
	var weightResponses []Network.MainData
	
	//Timestamp variable
	task_dist_timestamp := time.Now()
	
	//Variables used for channels
	var message_distributor_to_network Network.MainData
	var message_distributor_from_network Network.MainData
	var internal_message_distributor_to_task_manager Network.InternalMessage
	var internal_message_distributor_from_task_manager Network.InternalMessage
	
	fmt.Println("Just so Go doesnt complain about varible usage: ",message_distributor_to_network, message_distributor_from_network, internal_message_distributor_to_task_manager, internal_message_distributor_from_task_manager )
	
	//TEMP: Writing out starting timestamp
	fmt.Println("Timestamp: ",task_dist_timestamp)
	
	
	
	
	
	for {
		switch task_distributor_state{
		case waiting_for_task:
			
			select{
				case internal_message_distributor_from_task_manager := <-channel_from_task_manager:
					if internal_message_distributor_from_task_manager.Message_type == Network.ID_MSG_TYPE_DISTRIBUTOR_NEW_COMMAND {
						
						//Update current order variable
						currentOrder = internal_message_distributor_from_task_manager.Data
						//Printfs, REMOVE later
						fmt.Println("Received new order from Task Manager: ", internal_message_distributor_from_task_manager.Data)
						
						//Construct message to Elevator Controller (via Network) ID_MSG_TYPE_ELEVATOR_CONTROLLER_REQUEST_WEIGHTS
						message_distributor_to_network.Destination = Network.DESTINATION_BROADCAST
						message_distributor_to_network.Message_type = Network.ID_MSG_TYPE_ELEVATOR_CONTROLLER_REQUEST_WEIGHTS
						message_distributor_to_network.Data = [][]int{currentOrder}
						//Update timestamp and clear response array
						task_dist_timestamp = time.Now()
						weightResponses = nil
						//Send over network channel
						channel_to_network <- message_distributor_to_network
						
						
						task_distributor_state = waiting_for_response
					} else {
						//Received some other message
						fmt.Println("Received some other message")
						fmt.Println("MsgTypeID given: ", internal_message_distributor_from_task_manager.Message_type)
					}
					
				default:
					//No message ready for read. Do nothing

			
			}
			
			
		case waiting_for_response:
			select{
				case message_distributor_from_network := <- channel_from_network:
					weightResponses = append(weightResponses, message_distributor_from_network)
					
				default:
					if time.Since(task_dist_timestamp).Seconds() > Network.TIMEOUT_MESSAGE_RESPONSE{
						fmt.Println("Weight responses given: ",weightResponses)
						if len(weightResponses) <= 0 {
							fmt.Println("Did not receive any responses")
						} else {
							
							var responseChosen Network.MainData = weightResponses[0]
							for i := 0 ; i <len(weightResponses) ; i++{
								if weightResponses[i].Data[0][0] < responseChosen.Data[0][0] {
								responseChosen = weightResponses[i]
								}
							}
							
							message_distributor_to_network.Destination = responseChosen.Source
							message_distributor_to_network.Message_type = Network.ID_MSG_TYPE_ELEVATOR_HANDLER_DISTRIBUTED_ORDER
							message_distributor_to_network.Data = [][]int{currentOrder}
							
							channel_to_network <- message_distributor_to_network
																
						}
						task_distributor_state = waiting_for_task
						
					} else {
						//Do nothing
					}						
			}	
		default:
			task_distributor_state = waiting_for_task	
			
			fmt.Println("Time difference: ",time.Since(task_dist_timestamp).Seconds())
			fmt.Println("")
			
		}
	}
	
}
