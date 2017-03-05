
package Tester
import (
	"fmt"
	"time"
	"../../Network"
	"math/rand"
)


func Test_distributor_networkEmulator(send_chan chan Network.MainData, receive_chan chan Network.MainData) {
	fmt.Println("Tester for task distributor, Network Emulator")
	
	var network_message_receive Network.MainData
	var network_message_send Network.MainData
	
	fmt.Println("To allow execution: ", network_message_receive, network_message_send)
	
	for{
		select{
			case network_message_receive := <- receive_chan:
			if network_message_receive.Message_type == Network.ID_MSG_TYPE_ELEVATOR_CONTROLLER_REQUEST_WEIGHTS {
				//Build response message
				network_message_send.Source = "10"+"a"
				network_message_send.Destination = network_message_receive.Source
				network_message_send.Message_type = Network.ID_MSG_TYPE_DISTRIBUTOR_RESPONSE_WEIGHTS
				
				iterates := rand.Intn(6)
				
				for i := 0 ; i <= iterates ; i++ {
					network_message_send.Source +="a"
					randomWeight := []int{rand.Intn(254)}
					network_message_send.Data = [][]int{randomWeight}
					send_chan <- network_message_send
					time.Sleep(250*time.Millisecond)
				}
				
				
				
				
			} else if network_message_receive.Message_type == Network.ID_MSG_TYPE_ELEVATOR_HANDLER_DISTRIBUTED_ORDER {
				fmt.Println("Distributor distributed the following order: ",network_message_receive)
				fmt.Println("")
				fmt.Println("")
				fmt.Println("")
			} else {
				fmt.Println("Network emulator received an unknown message")
			}
			
			default:
			//Do nothing
			
			
		}
		
		
		
	}

}


func Test_distributor_taskmanagerEmulator(send_chan chan Network.InternalMessage, receive_chan chan Network.InternalMessage){
	fmt.Println("Tester for task distributor, task manager Emulator")
	
	var internal_message_test_to_distributor Network.InternalMessage
	for{
		
		internal_message_test_to_distributor.Message_type = Network.ID_MSG_TYPE_DISTRIBUTOR_NEW_COMMAND
		randomData := []int{rand.Intn(5),rand.Intn(2)}
		internal_message_test_to_distributor.Data = randomData
		time.Sleep(5*time.Second)
		
		send_chan <- internal_message_test_to_distributor
	}	
	
	
	
}