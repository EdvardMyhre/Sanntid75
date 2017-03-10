
package Tester
import (
	"fmt"
	"time"
	"../../Network"
	"math/rand"
)



func Test_pendingAndBackup_manager_buttonIntermediarySimulator(send_chan chan<- Network.Button){
	fmt.Println("Simulator Running: Button Intermediary for Pending Tasks and Backup manager")
	
	var message_output Network.Button
	
	for{
		random_number := rand.Intn(6)
		message_output.Floor = random_number
		time.Sleep(4*time.Second)
		send_chan <- message_output
	}
	
	
}