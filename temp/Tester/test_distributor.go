
package Tester
import (
	"fmt"
	"time"
	"../../Network"
	"math/rand"
)


func Test_distributor_networkEmulator(send_chan chan Network.MainData, receive_chan chan Network.MainData) {
	fmt.Println("Tester for task distributor, Network Emulator")
	

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