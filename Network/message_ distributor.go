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

func Message_distributor(port int) {
    fmt.Println("er i Message_distributor")
    var message MainData
    fmt.Println(port)
    message = Udp_listner(port)
    fmt.Printf("Source = %v, Destination = %v, Message_type = %v, Data = %v", message.Source, message.Destination, message.Message_type, message.Data)
    fmt.Println()

}






