package main

import "./types"

//import "./amanager"
//import "./driver"

import "fmt"

//import "time"

// func main() {
// 	//TEST OF AMANAGER
// 	driver.Init()

// 	//Initialization of variables
// 	elev_status_c := make(chan types.Status)
// 	elev_task_c := make(chan int)
// 	pmanager_task_c := make(chan types.Task)
// 	pmanager_status_c := make(chan types.Task)
// 	udp_rx_c := make(chan types.MainData)
// 	udp_tx_c := make(chan types.MainData)

// 	elev_task := 0
// 	//pmanager_status := types.Task{}
// 	udp_in := types.MainData{}
// 	udp_out := types.MainData{}

// 	l := 3
// 	w := 4
// 	var a [][]int
// 	a = make([][]int, l)
// 	for i := 0; i < l; i++ {
// 		a[i] = make([]int, w)
// 	}

// 	a[0][0] = 0
// 	a[0][1] = 3
// 	a[0][2] = 0
// 	a[0][3] = 255

// 	a[1][0] = 2
// 	a[1][1] = 2
// 	a[1][2] = 0
// 	a[1][3] = 255

// 	a[2][0] = 1
// 	a[2][1] = 3
// 	a[2][2] = 0
// 	a[2][3] = 255

// 	go amanager.AssignedTasksManager(elev_status_c, elev_task_c,
// 		pmanager_task_c, pmanager_status_c,
// 		udp_rx_c, udp_tx_c)

// 	for {
// 		select {
// 		case elev_task = <-elev_task_c:
// 			fmt.Println("MAIN: New elevator task:", elev_task)
// 		case pmanager_status := <-pmanager_status_c:
// 			fmt.Println("MAIN: New status to pmanager:", pmanager_status)
// 		case udp_in = <-udp_tx_c:
// 			if udp_in.Destination == "backup" && udp_in.Type == types.REQUEST_BACKUP {
// 				udp_out.Type = types.GIVE_BACKUP
// 				udp_out.Data = a
// 				select {
// 				case udp_rx_c <- udp_out:
// 				case <-time.After(time.Second * 10):
// 					fmt.Println("MAIN: Could not send backup")
// 				}

// 			}

// 		}
// 	}

// }
