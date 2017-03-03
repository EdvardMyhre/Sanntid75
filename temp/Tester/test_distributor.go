
package Tester
import (
	"fmt"
	"time"
	"../../Network"
)


func Test_distributor_networkEmulator(send_chan chan Network.MainData, receive_chan chan Network.MainData) {
	fmt.Println("Tester for task distributor, Network Emulator")
	

}


func Test_distributor_taskmanagerEmulator(send_chan chan Network.NewOrder, receive_chan chan Network.NewOrder){
	fmt.Println("Tester for task distributor, task manager Emulator")
	
	var message_test_to_distributor Network.NewOrder
	message_test_to_distributor.Message_type = 1
	for{
		time.Sleep(1*time.Second)
		fmt.Println("Tester function: Task Manager Emulator Ping")
		
		send_chan <- message_test_to_distributor
		
		message_test_to_distributor.Message_type += 1
	}	
	
	
	
}