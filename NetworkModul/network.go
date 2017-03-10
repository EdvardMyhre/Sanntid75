package network

import (
	"./network/bcast"
	"./network/localip"
	"./network/messageid"
	"./network/peers"
	"./network/structer"
	"flag"
	"fmt"
	"os"
	//"sort"
	"math/rand"
	//"time"
)

func Network_start(n_to_distri chan structer.MainData, n_to_p_task_manager chan structer.MainData, n_to_a_tasks_manager chan structer.MainData,
					distri_to_n chan structer.MainData, p_task_manager_to_n chan structer.MainData, a_task_manager_to_n chan structer.MainData) {

	//------------------------ Kanaler mellom moduler -------------------------------------

	//assigned_tasks_manager_to_network := make(chan structer.MainData)
	//network_to_assigned_tasks_manager := make(chan structer.MainData)

	//distributing_state_machine_to_network := make(chan structer.MainData)
	//network_to_distributing_state_machine := make(chan structer.MainData)

	//task_manager_to_network := make(chan structer.MainData)
	//network_to_task_manager := make(chan structer.MainData)

	var message_send structer.MainData

	//-----------------------  Lager kanal som har oversikt over hvem som er i livet --------------------
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	//-----------------------  Lager kanaler for å sende og receive meldinger --------------------
	message_sendCh := make(chan structer.MainData)
	message_receivedCh := make(chan structer.MainData)


//-----------------------  Finner lokal ip og deklarerer myBackupId --------------------
	id := find_localip()
	myBackupId := ""

	//-----------------------  Sjekker hvem som er online --------------------
	go peers.Transmitter(50018, id, peerTxEnable)
	go peers.Receiver(50018, peerUpdateCh)

	//-----------------------  Sender og mottar fra broadcast --------------------
	go bcast.Transmitter(40018, message_sendCh)
	go bcast.Receiver(40018, message_receivedCh)

	//-----------------------  Melding for å sjekke om man har en backup --------------------
	/*go func() {
		message := structer.MainData{}
		message.Source = id
		message.Destination = "broadcast"
		message.Message_type = messageid.ID_MSG_TYPE_MY_BACKUP
		row1 := []int{}
		row2 := []int{}
		message.Data = append(message.Data, row1)
		message.Data = append(message.Data, row2)

		for {
			message_sendCh <- message
			time.Sleep(1 * time.Second)
		}
	}()*/
	//-----------------------  Fordeler det som kommer fra broadcast og håndterer hvem som er i livet --------------------
	go func(){
		for {
			select{ 
			case d := <-distri_to_n:
				message_send = d

			case p := <-p_task_manager_to_n:
				message_send = p

			case a := <-a_task_manager_to_n:
				message_send = a
			}
			message_send.Source = id
			message_sendCh <- message_send
		}
	}()






	//-----------------------  Fordeler det som kommer fra broadcast og håndterer hvem som er i livet --------------------
	fmt.Println("Started")
	go func(){
		for {
			select {
			case p := <-peerUpdateCh:
				if myBackupId == "" {
					myBackupId = find_backup(id, p)
				}
				//fmt.Println("Min id er:       ", id)
				//fmt.Println("Min backupid er: ", myBackupId)
				/*fmt.Printf("Peer update:\n")
				fmt.Printf("  Peers:    %q\n", p.Peers)
				fmt.Printf("  Backup:   %q\n", p.Backup)
				fmt.Printf("  New:      %q\n", p.New)
				fmt.Printf("  Lost:     %q\n", p.Lost)*/

			case a := <-message_receivedCh:
				if a.Destination == id || a.Destination == "broadcast" {
					switch a.Message_type & 224 {
					case messageid.ID_MODULE_DISTRIBUTOR:
						n_to_distri <- a

					case messageid.ID_MODULE_TASK_MANAGER:
						n_to_p_task_manager <- a

					case messageid.ID_MODULE_ELEVATOR_CONTROLLER:
						n_to_a_tasks_manager <- a

					case messageid.ID_MODULE_NETWORK:
						backup(a, myBackupId, id, message_sendCh)
						//fmt.Println("melding fra message_receivedCh: ", a)

					}
				}
			}
		}
	}()
}






//--------------------- Finner din backup  -----------------------------
/*func find_backup(id string, p peers.PeerUpdate) string {
	index := sort.SearchStrings(p.Peers, id)
	myBackupId := ""
	for i := range p.Backup {
		if i == index && p.Backup[i] != id {
			myBackupId = p.Backup[i]
			break
		}
	}
	return myBackupId
}*/

func find_backup(id string, p peers.PeerUpdate) string {
	if len(p.Peers) > 1 {
		for {
			i := rand.Intn(len(p.Peers))
			//fmt.Println(i)
			myBackupId := p.Peers[i]
			if myBackupId != id{
				return myBackupId
			}
		}
	}
	return ""
}






//-----------------------  Finner lokalip --------------------
func find_localip() string {

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("%s-%d", localIP, os.Getpid())
	}
	return id
}

func backup(m structer.MainData, myBackupId string, id string, message_sendCh chan structer.MainData) {
	if (m.Destination == "broadcast") && ((m.Message_type & 31) == messageid.ID_MSG_TYPE_IS_MY_BACKUP_ALIVE){
		//fmt.Println("heieie")
		if m.Source == myBackupId{
		message := structer.MainData{}
		message.Source = id
		message.Destination = m.Source
		message.Message_type = messageid.ID_MSG_TYPE_IS_MY_BACKUP_ALIVE_TRUE
		row1 := []int{}
		row2 := []int{}
		message.Data = append(message.Data, row1)
		message.Data = append(message.Data, row2)
		message_sendCh <- message	
		fmt.Println("message:  ", message)
		}
	}

}