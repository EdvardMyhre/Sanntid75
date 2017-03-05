package Network


import (
	"fmt"
	"net"
	//"os"
	"strings"
	//"time"
)

var localIP string
var port int = 16569



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



//------------------------------ message_from_net ---------------------------------------

func Message_from_net() {
    var message MainData
    localIp, _ := LocalIP()
    for {
        message = Udp_listner(port)

        if message.Destination == localIp || message.Destination == "broadcast" {
            Message_to_channel(message)
        }
    }

    //fmt.Printf("Source = %v, Destination = %v, Message_type = %v, Data = %v", message.Source, message.Destination, message.Message_type, message.Data)
    //fmt.Println()

}


//----------------------------- Message_from_modul --------------------------------------

//får meldingen fra interne moduler og broadcaster. sjekker om det er backup 

func Message_from_modul(message MainData) {
    localIp, _ := LocalIP()
    message.Source = localIp

    switch message.Destination{
        case localIp:
            Message_to_channel(message)

        case "broadcast": 
            Udp_broadcast(message, port)

        case "backup":
            //fin riktig backup og legg det i destination

    }

}


//------------------------------- Message_to_channel----------------------------

func Message_to_channel(message MainData) {
/*message_to_distributor_chan chan MainData, 
    message_to_task_manager_chan chan MainData, message_to_elevator_controller_chan chan MainData*/
    switch (message.Message_type & 224){
        case ID_MODULE_NETWORK:
            fmt.Println("message blir sendt til Network")
        case ID_MODULE_DISTRIBUTOR:
            fmt.Println("message blir sendt til Distributor")
        case ID_MODULE_TASK_MANAGER:
            fmt.Println("message blir sendt til Task_manager")
        case ID_MODULE_ELEVATOR_CONTROLLER:
            fmt.Println("message blir sendt til Elevator_controller")
        }
}




//-------------------------------- Hvem er backup ---------------------



