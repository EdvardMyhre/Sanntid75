package main

import (
	//"fmt"
	"./Network"
	"time"
)

//var port int = 16569

func main() {
	//a, _ := Network.LocalIP()

	send_objekt := Network.MainData{}
	send_objekt.Source = ""
	send_objekt.Destination = "broadcast"
	send_objekt.Message_type = 160
	row1 := []int{1, 2, 3, 4, 52}
	row2 := []int{4, 5, 6, 564, 4}
	send_objekt.Data = append(send_objekt.Data, row1)
	send_objekt.Data = append(send_objekt.Data, row2)
   

   	//chan_send <- send_objekt
	for i := 0; i < 60; i++ {
		Network.Message_from_modul(send_objekt)
		time.Sleep(time.Second * 1)
	}

}
