//-----------------------  Heis 1  ----------------

package main

import (
	//"fmt"
	//"net"
	"./Network"
	//"time"
)

func main() {

	a, _ := Network.LocalIP()

	send_objekt := Network.MainData{}
	send_objekt.Source = a
	send_objekt.Destination = "heis 2"
	//b, err := json.Marshal(data)
	send_objekt.Message_type = 3
	row1 := []int{1, 2, 3, 4, 52}
	row2 := []int{4, 5, 6, 564, 4}
	send_objekt.Data = append(send_objekt.Data, row1)
	send_objekt.Data = append(send_objekt.Data, row2)

	Network.Udp_listner(10001)
	/*
		for i := 0; i < 6; i++ {
			go Network.Udp_broadcast(send_objekt, 10001)
			time.Sleep(time.Second * 1)
		}*/

}
