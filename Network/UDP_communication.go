package Network


import (
	"fmt"
	"net"
	//"os"
	//"strings"
	"time"
)



var source string
var destination string
//var portNr string = ":10001"
//var BcastIP string = "129.241.187.255"
var BcastIP string = "255.255.255.255"
//var BcastIP string = "10.22.79.255"
//var BcastIP string = "10.22.77.79"
//var BcastIP string = "10.22.76.255"
//var BcastIP string = "192.30.253.125"




//------------------------------ check_error ---------------------------------------

func check_error(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        //os.Exit(0)
    }
}





//------------------------------ Udp_listner ---------------------------------------



func Udp_listner(port int) MainData {
	ServerAddr,err := net.ResolveUDPAddr("udp",fmt.Sprintf(":%d",port)) 
    check_error(err)
 
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    check_error(err)
    fmt.Println("er i udp_listner")


    buf := make([]byte, 1024)

    for {
    	defer ServerConn.Close()
       	time.Sleep(time.Second * 1)
        fmt.Println("hei er i udp_listner og for")
    	n,_,err := ServerConn.ReadFromUDP(buf)
        receivedsend_objekt := string(buf[0:n])
        receivedsend_objektBytes := []byte(receivedsend_objekt)
        

        if err == nil {
            data := Json_to_struct_MainData(receivedsend_objektBytes)
            //fmt.Printf("Source = %v, Destination = %v, send_objekt_type = %v, Data = %v", data.Source, data.Destination, data.send_objekt_type, data.Data)
            //fmt.Println()
            return data
        }
    }
}





//------------------------------ Udp_broadcast ---------------------------------------


func Udp_broadcast(data MainData, port int) {

	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", BcastIP ,port))
	conn, _ := net.DialUDP("udp", nil ,addr)

    defer conn.Close()

    send := Struct_to_json_MainData(data)
    _,err1 := conn.Write(send)
    if err1 != nil {
        fmt.Println(err1)
    }
}










