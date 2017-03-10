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
	"sort"
	"time"
)

func Network_start(network_to_distributing_state_machine chan structer.MainData) {

	//------------------------ Kanaler mellom moduler -------------------------------------

	//assigned_tasks_manager_to_network := make(chan structer.MainData)
	//network_to_assigned_tasks_manager := make(chan structer.MainData)

	//distributing_state_machine_to_network := make(chan structer.MainData)
	//network_to_distributing_state_machine := make(chan structer.MainData)

	//task_manager_to_network := make(chan structer.MainData)
	//network_to_task_manager := make(chan structer.MainData)

	//-----------------------  Lager kanal som har oversikt over hvem som er i livet --------------------
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	//-----------------------  Lager kanaler for å sende og receive meldinger --------------------
	message_send := make(chan structer.MainData)
	message_received := make(chan structer.MainData)

	id := find_localip()
	myBackupId := ""

	//-----------------------  Sjekker hvem som er online --------------------
	go peers.Transmitter(50018, id, peerTxEnable)
	go peers.Receiver(50018, peerUpdateCh)

	//-----------------------  Sender og mottar fra broadcast --------------------

	/*go func() {
		select {
		case m := <-message_send:
			switch m.Destination {
			case "backup":
				fmt.Println("kom in i case")
				m.Destination = myBackupId
				//message_send <- m
				fmt.Println("kom ut av case:  ", m.Destination)
			}
		}
		bcast.Transmitter(40018, message_send)
	}()*/
	go bcast.Transmitter(40018, message_send)
	go bcast.Receiver(40018, message_received)

	//-----------------------  Lager et objekt som sendes ut --------------------
	go func() {
		message := structer.MainData{}
		message.Source = id
		message.Destination = "backup"
		message.Message_type = messageid.ID_MODULE_DISTRIBUTOR
		row1 := []int{1, 2, 3, 4, 52}
		row2 := []int{4, 5, 6, 564, 4}
		message.Data = append(message.Data, row1)
		message.Data = append(message.Data, row2)

		for {
			var i int
			i++
			message_send <- message
			time.Sleep(1 * time.Second)
		}
	}()

	//-----------------------  Printer ut hva som kommer fra broadcast og hvem som er i livet --------------------
	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			if myBackupId == "" {
				myBackupId = find_backup(id, p)
			}
			fmt.Println("Min id er:       ", id)
			fmt.Println("Min backupid er: ", myBackupId)
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  Backup:   %q\n", p.Backup)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-message_received:
			//fmt.Println(a)
			switch a.Message_type & 224 {

			case messageid.ID_MODULE_DISTRIBUTOR:
				//fmt.Println(a)
				network_to_distributing_state_machine <- a

			case messageid.ID_MODULE_TASK_MANAGER:
				//network_to_task_manager <- a
				fmt.Println("Sender til task_manager      ", a)

			case messageid.ID_MODULE_ELEVATOR_CONTROLLER:
				fmt.Println("Sender til elvator controller      ", a)
				//network_to_assigned_tasks_manager <- a
			}
		}
	}
}

//--------------------- Finner din backup  -----------------------------
func find_backup(id string, p peers.PeerUpdate) string {
	index := sort.SearchStrings(p.Peers, id)
	myBackupId := ""
	for i := range p.Backup {
		if i == index && p.Backup[i] != id {
			myBackupId = p.Backup[i]
			break
		}
	}
	return myBackupId
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
