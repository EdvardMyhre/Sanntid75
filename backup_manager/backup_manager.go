
package backup_manager
import (
	"fmt"
	"time"
	"../../types"
)



type struct_backup_element struct {
	BackupIP      string
	BackupData    [][]int
	
	
}

const (
	TEMPREQUEST = 1
	TEMPPUSH = 2
)



func Backup_manager(	channel_from_network 	<-chan types.MainData,		channel_to_network 			chan<- types.MainData,
																			channel_to_pending_manager 	chan<- types.MainData){
	//var message_maindata types.MainData
//CREATE MATRIX
	var backup_matrix []struct_backup_element
	var sendQueue_push []string
	var sendQueue_request []string
	var sendQueue_lostbackup []string

	for {
//BEHAVIOR FOR RECEIVING FROM NETWORK
		select {
			case network_message := <- channel_from_network:
				fmt.Println("Received something from network: ",network_message)
		//"REQUEST BACKUP" RECEIVED
				if network_message.Type == TEMPREQUEST {
					var request_command_already_exists bool
					request_command_already_exists = false
					for i := 0 ; i < len(sendQueue_request) ; i++ {
						if sendQueue_request[i] == network_message.Source {
							request_command_already_exists = true
							break
						}
					}
					if request_command_already_exists == false {
						sendQueue_request = append(sendQueue_request, network_message.Source)
						fmt.Println("Received NEW REQUEST command message. Added to queue. Current request queue: ", sendQueue_request)
					}
				

		//"PUSH BACKUP" RECEIVED		
				} else if network_message.Type == TEMPPUSH {
					var backup_already_exists bool
					var push_command_already_exists bool
					backup_already_exists = false
					push_command_already_exists = false
					//Check if a backup already exists
					for i := 0 ; i < len(backup_matrix) ; i++ {
						if backup_matrix[i].BackupIP == network_message.Source {
							backup_matrix[i].BackupData = network_message.Data
							backup_already_exists = true
							break
						}
					}
					//if not, create new element and append to matrix
					if backup_already_exists == false {
						var new_backup struct_backup_element
						new_backup.BackupIP = network_message.Source
						new_backup.BackupData = network_message.Data
						backup_matrix = append(backup_matrix, new_backup)
					}
					//Check if response command already exists
					for i:= 0 ; i < len(sendQueue_push) ; i++ {
						if sendQueue_push[i] == network_message.Source {
							push_command_already_exists = true
						}
					}
					
					//if not, create new command and append to list
					if push_command_already_exists == false {
						sendQueue_push = append(sendQueue_push, network_message.Source)
						fmt.Println("Received NEW push command message. Added to queue. Current queue: ", sendQueue_push)
					}
					
		//"BACKUP LOST" RECEIVED	
				} else if network_message.Type == types.BACKUP_LOST {
					//Check if task is already in queue
					var command_already_exists bool
					command_already_exists = false
					for i := 0 ; i< len(sendQueue_lostbackup) ; i++ {
						if sendQueue_lostbackup[i] == network_message.Source {
							command_already_exists = true
							fmt.Println("Received an already exsting backup lost message. Duplicate NOT added to queue. Current queue: ",sendQueue_lostbackup)
						}
					}
					if command_already_exists == false {
						sendQueue_lostbackup = append(sendQueue_lostbackup, network_message.Source)
						fmt.Println("Received new backup lost message. Current queue: ",sendQueue_lostbackup)
					}
				}
			
			default:
			//Do nothing
		}
		
//BEHAVIOR FOR SENDING "PUSH" RESPONSE TO NETWORK
		//If there's a an order in network sendQueue
		if len(sendQueue_push) > 0 {
			var push_index int
			var push_backup_exists bool
			push_backup_exists = false
			for i := 0 ; i<len(backup_matrix) ; i++ {
				if backup_matrix[i].BackupIP == sendQueue_push[0] {
					push_index = i
					push_backup_exists = true
					break
				}
			}
			
			var push_message types.MainData
			push_message.Destination = sendQueue_push[0]
			if push_backup_exists {
				push_message.Data = backup_matrix[push_index].BackupData
			}
			select{
				case channel_to_network <- push_message:
					fmt.Println("Sent push response back to: ",sendQueue_push[0])
					//Delete element in index 0
					sendQueue_push = append(sendQueue_push[1:])
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
			}
			
		}
//BEHAVIOR FOR SENDING "REQUEST" RESPONSE TO NETWORK		
		if len(sendQueue_request) > 0 {
			var request_index int
			var request_backup_exists bool
			request_backup_exists = false
			for i := 0 ; i<len(backup_matrix) ; i++ {
				if backup_matrix[i].BackupIP == sendQueue_request[0] {
					request_index = i
					request_backup_exists = true
					break
				}
			}
			
			var request_message types.MainData
			request_message.Destination = sendQueue_request[0]
			if request_backup_exists {
				request_message.Data = backup_matrix[request_index].BackupData
			}
			select{
				case channel_to_network <- request_message:
				fmt.Println("Sent request response back to: ",sendQueue_request[0])
				sendQueue_request = append(sendQueue_request[1:])
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
			}
			
			
		}
		
//BEHAVIOR FOR SENDING "BACKUP LOST" TO PENDING MANAGER
		if len(sendQueue_lostbackup) > 0 {
			var lost_backup_index int
			var lost_backup_exists bool
			lost_backup_exists = false
			
			for i := 0 ; i<len(backup_matrix) ; i++ {
				if backup_matrix[i].BackupIP == sendQueue_lostbackup[0] {
					lost_backup_index = i
					lost_backup_exists = true
					break
				}
			}
			
			if lost_backup_exists {
				var lostbackup_message types.MainData
				lostbackup_message.Data = backup_matrix[lost_backup_index].BackupData
				select {
					case channel_to_pending_manager <- lostbackup_message:
						fmt.Println("Sent lost backup matrix to pending. Lost backup IP: ", sendQueue_lostbackup[0])
						sendQueue_lostbackup = append(sendQueue_lostbackup[1:])
					case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
				}
			}
		}
	}
																				
}
