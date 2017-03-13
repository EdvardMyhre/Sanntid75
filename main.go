package main

import (
	"./NetworkModul"
	"./NetworkModul/network/structer"
	//"./NetworkModul/network/messageid"
	"fmt"
	//"time"
)

func main() {
	n_to_distri := make(chan structer.MainData)
	n_to_p_task_manager := make(chan structer.MainData)
	n_to_a_tasks_manager := make(chan structer.MainData)


	distri_to_n := make(chan structer.MainData)
	p_task_manager_to_n := make(chan structer.MainData)
	a_task_manager_to_n := make(chan structer.MainData)


	go network.Network_start(n_to_distri, n_to_p_task_manager, n_to_a_tasks_manager,
							distri_to_n, p_task_manager_to_n, a_task_manager_to_n)


	/*go func() {
		message := structer.MainData{}
		message.Source = ""
		message.Destination = "broadcast"
		message.Message_type = messageid.ID_MSG_TYPE_IS_MY_BACKUP_ALIVE
		row1 := []int{}
		row2 := []int{}
		message.Data = append(message.Data, row1)
		message.Data = append(message.Data, row2)

		for {
			distri_to_n <- message
			time.Sleep(1 * time.Second)
		}
	}()*/


	 for {
		select {
		case p := <-n_to_distri:
	 		fmt.Println("Sendt til dist:  ", p)
	 	case p := <-n_to_a_tasks_manager:
	 		fmt.Println("Sendt til a_task_manager:  ", p)
	 	case p := <-n_to_p_task_manager:
	 		fmt.Println("Sendt til p_task_manager:  ", p)
	 	}
	 }

}
