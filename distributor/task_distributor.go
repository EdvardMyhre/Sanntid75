
package main
import (
	"fmt"
	"time"
	"../Tester"
	//"../Network"
	"../types"
)

func main(){
	
	channel_distributor_to_network := make(chan types.MainData)
	channel_network_to_distributor := make(chan types.MainData)
	
	channel_task_manager_to_distributor := make(chan types.Task)
	channel_distributor_to_task_manager := make(chan types.Task)
	
	go task_distributor(channel_distributor_to_task_manager, channel_task_manager_to_distributor,
						channel_distributor_to_network, channel_network_to_distributor)
						
	
	go Tester.Test_distributor_networkEmulator(channel_network_to_distributor, channel_distributor_to_network)
	
	go Tester.Test_distributor_taskmanagerEmulator(channel_task_manager_to_distributor, channel_distributor_to_task_manager)
	
	//forever loop so the main function does not terminate
	for{}

}


func task_distributor(	channel_to_task_manager chan types.Task, 	channel_from_task_manager chan types.Task,
						channel_to_network chan types.MainData, 				channel_from_network chan types.MainData) {
	
	fmt.Println("task distributor, Connectivity test")
	
	//Constants for state machine
	const (
		waiting_for_newOrder = 1		//001
		confirm_order = 2
		request_weights = 4
		waiting_for_weightResponses = 8	//010
		distribute_order = 16
	)
	
	//State machine variable
	var task_distributor_state int = waiting_for_newOrder

	var currentOrder types.Task
	//Response array
	var weightResponses []types.MainData
	//Timestamp variable
	task_dist_timestamp := time.Now()
	//Variables used for channels
	var message_distributingOrder types.MainData
	fmt.Println("Just so Go doesnt complain about varible usage: ",message_distributingOrder)

	
	for {
		switch task_distributor_state{
		case waiting_for_newOrder:
			select{
				case inputOrder := <-channel_from_task_manager:
					
					currentOrder = inputOrder
					fmt.Println("Received new order from Task Manager: ", currentOrder)	
					task_distributor_state = confirm_order			
				default:
					//No message ready for read. Do nothing

			
			}
		case confirm_order:
			//Send a confirmation message back to pending manager. (To start timer)
			currentOrder.Finished = 0
			currentOrder.Assigned = 255
			fmt.Println("CURRENT ORDER WHEN SENDING CONFIRMATION: ",currentOrder)
			select{
				case channel_to_task_manager <- currentOrder:
					//Set up data for network
					message_distributingOrder.Destination = "broadcast"
					message_distributingOrder.Type = types.REQUEST_WEIGHT
					var networkOrder []int
						networkOrder = nil
						networkOrder = append(networkOrder, currentOrder.Type)
						networkOrder = append(networkOrder, currentOrder.Floor)
						networkOrder = append(networkOrder, 0)
						networkOrder = append(networkOrder, 0)
						//fmt.Println("										REQUEST WEIGHT NETWORKORDER: ",networkOrder)
						
					message_distributingOrder.Data = [][]int{networkOrder}
						//fmt.Println("										REQUEST WEIGHT DATA FIELD: ",message_distributingOrder.Data)
					task_distributor_state = request_weights
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
				
			}
		
		case request_weights:
			select{
				case channel_to_network <- message_distributingOrder:
					//Update timestamp and clear response array
					task_dist_timestamp = time.Now()
					weightResponses = nil
					
					task_distributor_state = waiting_for_weightResponses
					
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
			}
			
			
			
		case waiting_for_weightResponses:
			select{
				case message_distributingOrder := <- channel_from_network:
					if message_distributingOrder.Type == types.GIVE_WEIGHT {
						weightResponses = append(weightResponses, message_distributingOrder)
					}
				default:
					if time.Since(task_dist_timestamp) > types.TIMEOUT_NETWORK_MESSAGE_RESPONSE{
						task_distributor_state = distribute_order
					} else {
						//Do nothing
					}						
			}	
		case distribute_order:
			fmt.Println("Weight responses given: ",weightResponses)
				if len(weightResponses) <= 0 {
					fmt.Println("Did not receive any responses")
				} else {		
					var responseChosen types.MainData = weightResponses[0]
					for i := 0 ; i <len(weightResponses) ; i++{
						if weightResponses[i].Data[0][0] < responseChosen.Data[0][0] {
							responseChosen = weightResponses[i]
						}
					}
							
					message_distributingOrder.Destination = responseChosen.Source
					message_distributingOrder.Type = types.DISTRIBUTE_ORDER
					
					var networkOrder []int
						networkOrder = nil
						networkOrder = append(networkOrder, currentOrder.Type)
						networkOrder = append(networkOrder, currentOrder.Floor)
						networkOrder = append(networkOrder, 0)
						networkOrder = append(networkOrder, 0)
					
					
					message_distributingOrder.Data = [][]int{networkOrder}
							
					channel_to_network <- message_distributingOrder
					
					currentOrder.Finished = 1
					currentOrder.Assigned = 0
					channel_to_task_manager <- currentOrder
					
					
					task_distributor_state = waiting_for_newOrder
																
				}		
		
		default:
			task_distributor_state = waiting_for_newOrder	
			
			fmt.Println("Time difference: ",time.Since(task_dist_timestamp).Seconds())
			fmt.Println("")
			
		}
	}
	
}
