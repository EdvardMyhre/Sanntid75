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
	"time"
)

func Network_start(n_to_distri chan structer.MainData, n_to_p_task_manager chan structer.MainData, n_to_a_tasks_manager chan structer.MainData,
	distri_to_n chan structer.MainData, p_task_manager_to_n chan structer.MainData, a_task_manager_to_n chan structer.MainData) {

	//-----------------------  Finner lokal ip og deklarerer myBackupId --------------------
	id := find_localip()
	myBackupId := ""
	myBackupAlive := false
	backupFor := []string{}

	var message_send structer.MainData

	//-----------------------  Lager kanal som har oversikt over hvem som er i livet --------------------
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	//-----------------------  Lager kanaler for å sende og receive meldinger --------------------
	message_sendCh := make(chan structer.MainData)
	message_receivedCh := make(chan structer.MainData)

	//-----------------------  Sjekker hvem som er online --------------------
	go peers.Transmitter(50018, id, peerTxEnable)
	go peers.Receiver(50018, peerUpdateCh)

	//-----------------------  Sender og mottar fra broadcast --------------------
	go bcast.Transmitter(40018, message_sendCh)
	go bcast.Receiver(40018, message_receivedCh)

	//-----------------------  Får meldinger fra modul og sender til broadcast --------------------
	go func() {
		for {
			select {
			case d := <-distri_to_n:
				message_send = d

			case p := <-p_task_manager_to_n:
				message_send = p

			case a := <-a_task_manager_to_n:
				message_send = a
			}
			message_send.Source = id

			if message_send.Destination == "backup" {
				message_send.Destination = myBackupId
			}
			message_sendCh <- message_send
		}
	}()

	//-----------------------  Fordeler det som kommer fra broadcast og håndterer hvem som er i livet --------------------
	fmt.Println("Started")
	go func() {
		for {
			select {
			case p := <-peerUpdateCh:
				if myBackupId == "" {
					//for i := 0; i < 5; i++ {
						//send_message_is_my_backup_alive(id, message_sendCh)
					//}
					//time.Sleep(200 * time.Millisecond)
					//if myBackupAlive == false {
						//fmt.Println("leter etter backup da.....................")
						find_backup(id, p, &myBackupAlive, message_sendCh, &myBackupId)

					//}
				}

				//my_backup_is_gone(&myBackupAlive, backupFor, p, &myBackupId)
				time.Sleep(1000 * time.Millisecond)
				fmt.Println("Min id er:              ", id)
				fmt.Println("Min backupid er:        ", myBackupId)
				fmt.Println("Min myBackupAlive er:   ", myBackupAlive)
				//fmt.Println("backupFor:              ", backupFor)
				fmt.Println("")

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
						//fmt.Println("motat melding:   ", a)
						message_receive_backup_alive(id, a, backupFor, message_sendCh)
						my_backup_is_alive(id, &myBackupAlive, a)
						backup_for(id, a, &backupFor)
						fmt.Println("backupFor:              ", backupFor)
					}
				}
			}
		}
	}()
}

//--------------------- Finner din backup  -----------------------------
func find_backup(id string, p peers.PeerUpdate, myBackupAlive *bool, message_sendCh chan structer.MainData, myBackupId *string) {
	fmt.Println("inne i find_backup xxxx---------------------")

	if len(p.Peers) > 1 {
		for {
			i := rand.Intn(len(p.Peers))
			if p.Peers[i] != id {
				*myBackupAlive = true
				*myBackupId = p.Peers[i]

				message := structer.MainData{}
				message.Source = id
				message.Destination = *myBackupId
				message.Message_type = messageid.ID_MSG_TYPE_YOU_ARE_MY_BACKUP
				row1 := []int{}
				row2 := []int{}
				message.Data = append(message.Data, row1)
				message.Data = append(message.Data, row2)
				message_sendCh <- message
				//fmt.Println("find_backup: ", message)
				//time.Sleep(50 * time.Millisecond)
				return
			}
		}
	}
}

//--------------------------- Legger hvem du er backupfor i en liste ----------------------
func backup_for(id string, a structer.MainData, backupFor *[]string) {
	if (a.Destination == id) && ((a.Message_type & 31) == messageid.ID_MSG_TYPE_YOU_ARE_MY_BACKUP) {
		*backupFor = append(*backupFor, a.Source)
		//fmt.Println("oooooooooooooooooooooooooooooooooooooooooooooooo")
		//fmt.Println("backup_for:    ", *backupFor)
	}
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

//-----------------------  Melding for å sjekke om man har en backup --------------------
func send_message_is_my_backup_alive(id string, message_sendCh chan structer.MainData) {
	message := structer.MainData{}
	message.Source = id
	message.Destination = "broadcast"
	message.Message_type = messageid.ID_MSG_TYPE_IS_MY_BACKUP_ALIVE
	row1 := []int{}
	row2 := []int{}
	message.Data = append(message.Data, row1)
	message.Data = append(message.Data, row2)
	message_sendCh <- message
	fmt.Println("-------------------------------------------------")
}

//--------------------------- Min Backup lever ------------------------------
func message_receive_backup_alive(id string, m structer.MainData, backupFor []string, message_sendCh chan structer.MainData) {
	if (m.Destination == "broadcast") && ((m.Message_type & 31) == messageid.ID_MSG_TYPE_IS_MY_BACKUP_ALIVE) {
		for i := range backupFor {
			if backupFor[i] == m.Source {
				message := structer.MainData{}
				message.Source = id
				message.Destination = m.Source
				message.Message_type = messageid.ID_MSG_TYPE_IS_MY_BACKUP_ALIVE_TRUE
				row1 := []int{}
				row2 := []int{}
				message.Data = append(message.Data, row1)
				message.Data = append(message.Data, row2)
				message_sendCh <- message
				fmt.Println("message_receive_backup_alive:    ", message)

			}
		}
	}
}

func my_backup_is_alive(id string, myBackupAlive *bool, m structer.MainData) {
	if (m.Destination == id) && ((m.Message_type & 31) == messageid.ID_MSG_TYPE_IS_MY_BACKUP_ALIVE_TRUE) {
		*myBackupAlive = true
		fmt.Println("my_backup_is_alive:    ", *myBackupAlive)
	}

}

func my_backup_is_gone(myBackupAlive *bool, backupFor []string, p peers.PeerUpdate, myBackupId *string) {
	for i := range backupFor {
		for j := range p.Lost {
			if backupFor[i] == p.Lost[j] {
				*myBackupAlive = false
				*myBackupId = ""

			}
		}
	}
}

// Hvis den vi er backup for dør, send melding til pending_task_mangager
