package main

import "./types"

import "./amanager"
import "./driver"
import "./elevator"

import "fmt"
import "time"

func main() {
	//TEST OF AMANAGER
	driver.Init()

	//Initialization of variables
	elev_status_c := make(chan types.Status)
	elev_task_c := make(chan int)
	pmanager_task_c := make(chan types.Task)
	pmanager_status_c := make(chan types.Task)
	udp_rx_c := make(chan types.MainData)
	udp_tx_c := make(chan types.MainData)

	//pmanager_status := types.Task{}
	udp_in := types.MainData{}
	udp_out := types.MainData{}

	l := 1
	w := 4
	var a [][]int
	a = make([][]int, l)
	for i := 0; i < l; i++ {
		a[i] = make([]int, w)
	}

	//Hall up from floor 0
	a[0][0] = types.BTN_TYPE_UP
	a[0][1] = 0
	a[0][2] = 0
	a[0][3] = 255

	go elevator.Controller(elev_task_c, elev_status_c)
	go amanager.AssignedTasksManager(elev_status_c, elev_task_c,
		pmanager_task_c, pmanager_status_c,
		udp_rx_c, udp_tx_c)
	go elevator.ButtonPoller(pmanager_task_c)
	tries_request := 0
	tries_push := 0
	for {
		select {
		case pmanager_status := <-pmanager_status_c:
			fmt.Println("AMANAGER: New status to pmanager:", pmanager_status)
		case udp_in = <-udp_tx_c:
			if udp_in.Destination == "backup" && udp_in.Type == types.REQUEST_BACKUP {
				if tries_request > 3 {
					tries_request = 0
					udp_out.Type = types.GIVE_BACKUP
					udp_out.Data = a
					select {
					case udp_rx_c <- udp_out:
					case <-time.After(time.Second * 10):
						fmt.Println("MAIN: Could not give backup")
					}
				}
				tries_request++

			} else if udp_in.Destination == "backup" && udp_in.Type == types.PUSH_BACKUP {
				if tries_push > 3 {
					tries_push = 0
					udp_out.Type = types.ACK_BACKUP
					udp_out.Data = udp_in.Data
					select {
					case udp_rx_c <- udp_out:
					case <-time.After(time.Second * 10):
						fmt.Println("MAIN: Could not acknowledge backup")
					}
					fmt.Println("AMANAGER: pushed backup:")
					fmt.Println(udp_in.Data)
				}
				tries_push++

			} else if udp_in.Destination == "broadcast" && udp_in.Type == types.SET_LIGHT {
				//fmt.Println("AMANAGER: pushed light at button type:", udp_in.Data[0][0], "and floor:", udp_in.Data[0][1], "with finished:", udp_in.Data[0][2])
			} else if udp_in.Type == types.GIVE_WEIGHT {
				fmt.Println("AMANAGER: gave weight")
			}

		case <-time.After(time.Millisecond * 10):
		}
	}

}
