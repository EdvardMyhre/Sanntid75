
package main
import (
	"fmt"
	"time"
	"./Tester"
)

func main(){
	go task_distributor()
	//forever loop so the main function does not terminate
	go Tester.Test_distributor()
	for{}

}


func task_distributor() {
	fmt.Println("task distributor, Connectivity test")
	const (
		waiting_for_task = 1		//001
		waiting_for_response = 2	//010
		send_choice = 5				//100
									//111
	)
	var task_distributor_state int = waiting_for_task
	task_dist_timestamp := time.Now()
	fmt.Println("Timestamp: ",task_dist_timestamp)
	
	
	for {
		switch task_distributor_state{
		case waiting_for_task:
			//Waits for new task from task_manager
			//fmt.Println("Waiting for task")
			task_distributor_state = waiting_for_response
		case waiting_for_response:
			//Waits for other nodes to respond
			//fmt.Println("Waiting for response")
			task_distributor_state = 0
			
			
		case send_choice:
			//Sends message containing the node chosen for task
			
			
			
		default:
			task_distributor_state = waiting_for_task	
			fmt.Println("Default state")
			time.Sleep(2*time.Second)
			fmt.Println("Time difference: ",time.Since(task_dist_timestamp).Seconds())
			
		}
	}
	
}
