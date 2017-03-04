package Network


import (
	"fmt"
	"net"
	//"os"
	"strings"
	//"time"
)

var localIP string



//---------------------------- get local ip -------------------------------------

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
}



//------------------------------ message_distributor ---------------------------------------

func Message_from_net(port int) {
    var message MainData
    message = Udp_listner(port)
    switch (message.Message_type & 224){
    case ID_MODULE_NETWORK:
        fmt.Println("vi er i case 1")
    case ID_MODULE_DISTRIBUTOR:
        fmt.Println("vi er i case 2")
    case ID_MODULE_TASK_MANAGER:
        fmt.Println("vi er i case 3")
    case ID_MODULE_ELEVATOR_CONTROLLER:
        fmt.Println("vi er i case 4")
    }

    //fmt.Printf("Source = %v, Destination = %v, Message_type = %v, Data = %v", message.Source, message.Destination, message.Message_type, message.Data)
    //fmt.Println()

}

func Message_to_net() {
    var message MainData
    message = Udp_listner(port)
    switch (message.Message_type & 224){
    case ID_MODULE_NETWORK:
        fmt.Println("vi er i case 1")
    case ID_MODULE_DISTRIBUTOR:
        fmt.Println("vi er i case 2")
    case ID_MODULE_TASK_MANAGER:
        fmt.Println("vi er i case 3")
    case ID_MODULE_ELEVATOR_CONTROLLER:
        fmt.Println("vi er i case 4")
    }
}






