

//-----------------------  Heis 1  ----------------


package main 
 
import (
    //"fmt"
    //"net"
    "time"
    "./Network"
)




 
func main() {

    send_objekt := Network.MainData{}
    send_objekt.Source = "heis 1"
    send_objekt.Destination = "heis 2"
    //b, err := json.Marshal(data)
    send_objekt.Message_type = 3
    row1 := []int{1,2,3,4,52}
    row2 := []int{4,5,6,564,4}
    send_objekt.Data = append(send_objekt.Data,row1)
    send_objekt.Data = append(send_objekt.Data,row2)
    

    //time.Sleep(time.Second * 1)
    go Network.Udp_listner(10001)

    //fmt.Println("er her 1")

    for i := 0; i < 6; i++ {
        Network.Udp_broadcast(send_objekt, 10001)
        time.Sleep(time.Second * 1)
    }

    

   //time.Sleep(time.Second * 1)           //Kan brukes i en uenderlig forløkke


}