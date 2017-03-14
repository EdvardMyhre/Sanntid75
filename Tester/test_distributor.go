
package Tester
import (
	"fmt"
	"time"
	//"../../Network"
	"../types"
	"math/rand"
)


func Test_distributor_networkEmulator(send_chan chan types.MainData, receive_chan chan types.MainData) {
	
	var network_message_receive types.MainData
	var network_message_send types.MainData
	
	fmt.Println("Emulator Running: Network Module for Distributor testing: ", network_message_receive, network_message_send)
	
	for{
		select{
			case network_message_receive := <- receive_chan:
			if network_message_receive.Type == types.REQUEST_WEIGHT {
				//Build response message
				fmt.Println("Message received from distributor: ",network_message_receive.Data)
				network_message_send.Source = "10"+"a"
				network_message_send.Destination = network_message_receive.Source
				network_message_send.Type = types.GIVE_WEIGHT
				
				iterates := rand.Intn(6)
				
				for i := 0 ; i <= iterates ; i++ {
					network_message_send.Source +="a"
					randomWeight := []int{rand.Intn(254)}
					network_message_send.Data = [][]int{randomWeight}
					select{
						case send_chan <- network_message_send:
						case <-time.After(250*time.Millisecond):
					}
					
				}
				
				
				
				
			} else if network_message_receive.Type == types.DISTRIBUTE_ORDER {
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


func Test_distributor_taskmanagerEmulator(send_chan chan types.Task, receive_chan chan types.Task){
	fmt.Println("Emulator Running: Task Manager for Distributor testing")
	
	var internal_message_test_to_distributor types.Task
	var time_duration time.Time
	time_duration = time.Now()
	for{
		if(time.Since(time_duration) > 5*time.Second) {
			randomFloor := rand.Intn(3)
			randomOrder := rand.Intn(2)
			internal_message_test_to_distributor.Floor = randomFloor
			internal_message_test_to_distributor.Type = randomOrder
			select{
				case send_chan <- internal_message_test_to_distributor:
					time_duration = time.Now()
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
			}
		}

		select{
			case msg := <- receive_chan:
				fmt.Println("PENDING MANAGER RECEIVED MESSAGE FROM DISTRIBUTOR: ",msg)
			default:
		}
		
		
	}	
	
	
	
}