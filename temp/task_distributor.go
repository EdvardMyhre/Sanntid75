
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


func task_distributor(	channel_to_task_manager chan Network.InternalMessage, 	channel_from_task_manager chan Network.InternalMessage,
						channel_to_network chan Network.MainData, 				channel_from_network chan Network.MainData) {
	
	fmt.Println("task distributor, Connectivity test")
	
	//Constants for state machine
	const (
		waiting_for_newOrder = 1		//001
		waiting_for_weightResponses = 2	//010
	)
	
	//State machine variable
	var task_distributor_state int = waiting_for_newOrder
	//currentOrder = {floor, direction}
	var currentOrder []int
	//Response array
	var weightResponses []Network.MainData
	//Timestamp variable
	task_dist_timestamp := time.Now()
	//Variables used for channels
	var message_distributingOrder Network.MainData
	var message_orderDetails Network.InternalMessage
	fmt.Println("Just so Go doesnt complain about varible usage: ",message_distributingOrder, message_orderDetails )

	
	for {
		switch task_distributor_state{
		case waiting_for_newOrder:
			
			select{
				case message_orderDetails := <-channel_from_task_manager:
					if message_orderDetails.Message_type == Network.MESSAGE_TYPE_DISTRIBUTE_NEWORDER {
						
						//Update current order variable
						currentOrder = message_orderDetails.Data
						//Printfs, REMOVE later
						fmt.Println("Received new order from Task Manager: ", message_orderDetails.Data)
						
						//Construct message to Elevator Controller (via Network) ID_MSG_TYPE_ELEVATOR_CONTROLLER_REQUEST_WEIGHTS
						message_distributingOrder.Destination = Network.DESTINATION_BROADCAST
						message_distributingOrder.Message_type = Network.MESSAGE_TYPE_REQUEST_WEIGHT
						message_distributingOrder.Data = [][]int{currentOrder}
						//Update timestamp and clear response array
						task_dist_timestamp = time.Now()
						weightResponses = nil
						//Send over network channel
						channel_to_network <- message_distributingOrder
						
						
						task_distributor_state = waiting_for_weightResponses
					} else {
						//Received some other message
						fmt.Println("Received some other message")
						fmt.Println("Message Type received: ", message_orderDetails.Message_type)
					}
					
				default:
					//No message ready for read. Do nothing

			
			}
			
			
		case waiting_for_weightResponses:
			select{
				case message_distributingOrder := <- channel_from_network:
					weightResponses = append(weightResponses, message_distributingOrder)
					
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
							
							message_distributingOrder.Destination = responseChosen.Source
							message_distributingOrder.Message_type = Network.MESSAGE_TYPE_DISTRIBUTE_ORDER
							message_distributingOrder.Data = [][]int{currentOrder}
							
							channel_to_network <- message_distributingOrder
																
						}
						task_distributor_state = waiting_for_newOrder
						
					} else {
						//Do nothing
					}						
			}	
		default:
			task_distributor_state = waiting_for_newOrder	
			
			fmt.Println("Time difference: ",time.Since(task_dist_timestamp).Seconds())
			fmt.Println("")
			
		}
	}
	
}
