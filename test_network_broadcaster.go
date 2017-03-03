package main

import (
	//"fmt"
	"./Network"
	"time"
)

var port int = 16569

func main() {
	a, _ := Network.LocalIP()

	send_objekt := Network.MainData{}
	send_objekt.Source = a
	send_objekt.Destination = "heis 1"
	send_objekt.Message_type = 3
	row1 := []int{1, 2, 3, 4, 52}
	row2 := []int{4, 5, 6, 564, 4}
	send_objekt.Data = append(send_objekt.Data, row1)
	send_objekt.Data = append(send_objekt.Data, row2)

	//go Network.Udp_listner(port)

	for i := 0; i < 60; i++ {
		Network.Udp_broadcast(send_objekt, port)
		time.Sleep(time.Second * 1)
	}

}
