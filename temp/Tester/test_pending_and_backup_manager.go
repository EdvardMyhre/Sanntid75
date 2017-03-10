
package Tester
import (
	"fmt"
	"time"
	"../../Network"
	//"../../types"
	"math/rand"
)



func Test_pendingAndBackup_manager_buttonIntermediarySimulator(send_chan chan<- Network.Button){
	fmt.Println("Simulator Running: Button Intermediary for Pending Tasks and Backup manager")
	
	var message_output Network.Button
	
	for{
		random_number := rand.Intn(5)
		random_command := rand.Intn(3)
		message_output.Floor = random_number
		message_output.Button_type = random_command
		time.Sleep(2*time.Second)
		send_chan <- message_output
	}
	
	
}


func Test_pendingandBackup_manager_assignedSimulator(send_chan chan<- Network.Button, rec_chan <-chan Network.Button){
	//Send add/remove input to pending tasks manager
	var chan_message Network.Button
	
	for{
		random_number := rand.Intn(5)
		random_command := rand.Intn(3)
		random_addremove := rand.Intn(2)
		chan_message.Floor = random_number
		chan_message.Button_type = random_command
		chan_message.Add = random_addremove
		time.Sleep(4*time.Second)
		send_chan <- chan_message
	}
	
	
}