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
//var localIP string




//------------------------------ check_error ---------------------------------------

func check_error(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        //os.Exit(0)
    }
}




//------------------------------ LocalIP ---------------------------------------
/*
func LocalIP() (string, error) {
	if localIP == "" {
		conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
		if err != nil {
			return "", err
		}
		defer conn.Close()
		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	}
	return localIP, nil
}*/







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
            data := Udp_json_to_struct(receivedMessageBytes)
            fmt.Printf("Source = %v, Destination = %v, Message_type = %v, Data = %v", data.Source, data.Destination, data.Message_type, data.Data)
            fmt.Println()
        }
    }
}

/*func Udp_listner(port int) {
	ServerAddr,err := net.ResolveUDPAddr("udp",fmt.Sprintf(":%d",port)) 
    check_error(err)
 
 
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    check_error(err)
    fmt.Println(&ServerConn)

    //lastMessage := ""


    //ServerConn.SetReadDeadline(time.Now().Add(time.Second))

    buf := make([]byte, 1024)

    for {
    	defer ServerConn.Close()
       	time.Sleep(time.Second * 1)
    	n,_,err := ServerConn.ReadFromUDP(buf)
    	//fmt.Println(buf)
        receivedMessage := string(buf[0:n])
        receivedMessageBytes := []byte(receivedMessage)
        //check_error(err)

        if err == nil && lastMessage != receivedMessage {
        	lastMessage = receivedMessage
            data := Udp_json_to_struct(receivedMessageBytes)
            fmt.Printf("Source = %v, Destination = %v, Message_type = %v, Data = %v", data.Source, data.Destination, data.Message_type, data.Data)
            fmt.Println()

        }
       //ServerConn.SetReadDeadline(time.Now().Add(time.Second))

        if err == nil {
            data := Udp_json_to_struct(receivedMessageBytes)
            fmt.Printf("Source = %v, Destination = %v, Message_type = %v, Data = %v", data.Source, data.Destination, data.Message_type, data.Data)
            fmt.Println()

        }

    }
	
}*/






//------------------------------ Udp_broadcast ---------------------------------------


func Udp_broadcast(data MainData, port int) {

	//conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("255.255.255.255:%d", port))
	conn, _ := net.DialUDP("udp", nil ,addr)



	/*LocalAddr,err := net.ResolveUDPAddr("udp", port)
	check_error(err)

	Conn, err := net.DialUDP("udp", nil ,LocalAddr)
    check_error(err)*/

    defer conn.Close()

    send := Udp_struct_to_json(data)
    _,err1 := conn.Write(send)
    if err1 != nil {
        fmt.Println(err1)
    }
}


