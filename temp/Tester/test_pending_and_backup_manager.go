
package Tester
import (
	"fmt"
	"time"
	//"../../Network"
	"../../types"
	"math/rand"
)



func Test_pendingAndBackup_manager_buttonIntermediarySimulator(send_chan chan<- types.Button){
	fmt.Println("Simulator Running: Button Intermediary for Pending Tasks and Backup manager")
	
	var message_output types.Button
	
	for{
		random_number := rand.Intn(5)
		random_command := rand.Intn(3)
		message_output.Floor = random_number
		message_output.Type = random_command
		time.Sleep(2*time.Second)
		
		select {
			case send_chan <- message_output:
			
			case <-time.After(1*time.Second):
				fmt.Println("BUTTON INTERMEDIARY SEND FAIL")
		}


	}
	
	
}


func Test_pendingandBackup_manager_assignedSimulator(send_chan chan<- types.Task, rec_chan <-chan types.Task){
	//Send add/remove input to pending tasks manager
	var chan_message types.Task
	
	for{
		random_number := rand.Intn(5)
		random_command := rand.Intn(3)
		random_addremove := rand.Intn(2)
		chan_message.Floor = random_number
		chan_message.Type = random_command
		chan_message.Assigned = random_addremove
		time.Sleep(4*time.Second)
		select {
			case send_chan <- chan_message:
			case <-time.After(1*time.Second):
				fmt.Println("ASSIGNED SIMULATOR SEND FAIL")
		}
		
		select{
			case chan_message := <- rec_chan:
				fmt.Println("Received assigned task from pending manager: ",chan_message)
			default:
		}
		
	}
	
	
}

func Test_pendingandBackup_manager_distributorSimulator(send_chan chan<- types.Task, rec_chan <-chan types.Task){
	//Send periodic status updates
	var chan_message types.Task
	for{
		random_distributorstate := rand.Intn(2)
		chan_message.Finished = random_distributorstate
		time.Sleep(1*time.Second)
		send_chan <- chan_message
		
	}
	
	
	
}