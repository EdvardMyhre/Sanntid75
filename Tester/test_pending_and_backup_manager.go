
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
			
			case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
				//fmt.Println("BUTTON INTERMEDIARY SEND FAIL")
		}

		time.Sleep(10*time.Millisecond)
	}
	
	
}


func Test_pendingandBackup_manager_assignedSimulator(send_chan chan<- types.Task, rec_chan <-chan types.Task){
	//Send add/remove input to pending tasks manager
	var chan_message types.Task
	var time_duration time.Time
	
	time_duration = time.Now()
	for{
		if time.Since(time_duration) > 4*time.Second {
			
			random_number := rand.Intn(5)
			random_command := rand.Intn(3)
			random_addremove := rand.Intn(2)
			chan_message.Floor = random_number
			chan_message.Type = random_command
			chan_message.Assigned = random_addremove
			
			select {
				case send_chan <- chan_message:
					time_duration = time.Now()
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
					//fmt.Println("ASSIGNED SIMULATOR SEND FAIL")
				}	
		}
		select{
			case chan_message2 := <- rec_chan:
				//fmt.Println("Received assigned task from pending manager: ",chan_message)
				chan_message = chan_message2
			default:
		}
		time.Sleep(10*time.Millisecond)
	}
	
	
}

func Test_pendingandBackup_manager_distributorSimulator(send_chan chan<- types.Task, rec_chan <-chan types.Task){
	//Send periodic status updates
	var chan_message types.Task
	var time_duration time.Time
	
	time_duration = time.Now()
	
	for{
		if time.Since(time_duration) > 1*time.Second {
			random_distributorstate := rand.Intn(2)
			chan_message.Finished = random_distributorstate
			chan_message.Finished = 255
			//time.Sleep(1*time.Second)
			select{
				case send_chan <- chan_message:
					time_duration = time.Now()
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
					//fmt.Println("Distributor send FAILED<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			}
		}
		
		select{
			case chan_message2 := <- rec_chan:
				//fmt.Println("Task distributor received message; ",chan_message,"                                                                  DISTRIBUTOR RECEIVED MESSAGE")
				chan_message = chan_message2
			default:
		}
		
		time.Sleep(10*time.Millisecond)	
	}
	
	
	
}

func Test_pendingandBackup_manager_networkSimulator(send_chan chan<- types.MainData, rec_chan <-chan types.MainData) {
	fmt.Println("PENDING AND BACKUP MANAGER: NETWORK SIMULATOR STARTUP")
	var chan_message types.MainData
	var time_duration time.Time
	var sourceID string
	
	time_duration = time.Now()
	
	for{
		if time.Since(time_duration) > 2*time.Second {
			//network_message_send.Source = "10"+"a"
			//Create fake sourceID
			sourceID = "10"
			iterates := rand.Intn(6)		
			for i := 0 ; i <= iterates ; i++ {
				sourceID +="a"
			}
			chan_message.Source = sourceID
			chan_message.Type = types.BACKUP_LOST
			select{
				case send_chan <- chan_message:
					time_duration = time.Now()
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
			}
		
		}
	
		select{
			case chan_message2 := <- rec_chan:
				fmt.Println("Network received message from Backup Manager ", chan_message2)
			default:
				//Do nothing
		}
		time.Sleep(10*time.Millisecond)
	}
	
	
	
	
}