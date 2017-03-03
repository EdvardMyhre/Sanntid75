package Network


import (
	"fmt"
	"net"
	//"os"
	"strings"
	"time"
)



var source string
var destination string
//var portNr string = ":10001"
//var BcastIP string = "129.241.187.255"
//var BcastIP string = "255.255.255.255"
//var BcastIP string = "10.22.79.255"
//var BcastIP string = "10.22.77.79"
var BcastIP string = "10.22.77.255"




//------------------------------ check_error ---------------------------------------

func check_error(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        //os.Exit(0)
    }
}






//------------------------------ Udp_listner ---------------------------------------



func Udp_listner(port int) {
	ServerAddr,err := net.ResolveUDPAddr("udp",fmt.Sprintf(":%d",port)) 
    check_error(err)
 
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    check_error(err)


    buf := make([]byte, 1024)

    for {
    	defer ServerConn.Close()
       	time.Sleep(time.Second * 1)
    	n,_,err := ServerConn.ReadFromUDP(buf)
        receivedMessage := string(buf[0:n])
        receivedMessageBytes := []byte(receivedMessage)

        if err == nil {
            data := Json_to_struct_MainData(receivedMessageBytes)
            fmt.Printf("Source = %v, Destination = %v, Message_type = %v, Data = %v", data.Source, data.Destination, data.Message_type, data.Data)
            fmt.Println()
        }
    }
}





//------------------------------ Udp_broadcast ---------------------------------------


func Udp_broadcast(data MainData, port int) {

	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", BcastIP ,port))
	conn, _ := net.DialUDP("udp", nil ,addr)

    //fmt.Println(&conn)


    defer conn.Close()

    send := Struct_to_json_MainData(data)
    _,err1 := conn.Write(send)
    if err1 != nil {
        fmt.Println(err1)
    }
}






//---------------------------- get local ip -------------------------------------

var localIP string

func LocalIP() (string, error) {
    if localIP == "" {
        conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
        if err != nil {
            return "", err
        }
        defer conn.Close()
        localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
    }
    fmt.Println(localIP)
    return localIP, nil
}



