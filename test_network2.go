package main
 
import (
    //"fmt"
    "./Network"
    "time"
)




 
func main() {

    send_objekt := Network.MainData{}
    send_objekt.Source = "heis 2"
    send_objekt.Destination = "heis 1"
    send_objekt.Message_type = 3
    row1 := []int{1,2,3,4,52}
    row2 := []int{4,5,6,564,4}
    send_objekt.Data = append(send_objekt.Data,row1)
    send_objekt.Data = append(send_objekt.Data,row2)


    
    go Network.Udp_listner(10001)
    

    for i := 0; i < 6; i++ {
        Network.Udp_broadcast(send_objekt, 10001)
        time.Sleep(time.Second * 1)
    }

    

}