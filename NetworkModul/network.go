package network

import (
	"../types"
	"./network/bcast"
	"./network/localip"
	"./network/peers"
	"flag"
	"fmt"
	//"os"
	//"sort"
	"math/rand"
	"time"
)

func Network_start(n_to_distri chan types.MainData, n_to_p_task_manager chan types.MainData, n_to_a_tasks_manager chan types.MainData,
	distri_to_n chan types.MainData, p_task_manager_to_n chan types.MainData, a_task_manager_to_n chan types.MainData, n_to_a_tasks_manager2 chan types.MainData) {

	//-----------------------  Finner lokal ip og deklarerer myBackupId --------------------
	id := find_localip()
	myBackupId := ""
	myBackupAlive := false
	backupFor := []string{}

	var message_send types.MainData

	//-----------------------  Lager kanal som har oversikt over hvem som er i livet --------------------
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	//-----------------------  Lager kanaler for å sende og receive meldinger --------------------
	message_sendCh := make(chan types.MainData)
	message_receivedCh := make(chan types.MainData)

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
				//message_send.Destination = "broadcast"
			}
			message_sendCh <- message_send
			time.Sleep(types.PAUSE_NET_LISTNER)
		}
	}()

	//-----------------------  Fordeler det som kommer fra broadcast og håndterer hvem som er i livet --------------------
	go func() {
		for {
			select {
			case p := <-peerUpdateCh:
				if myBackupId == "" {
					for i := 0; i < 5; i++ {
						send_message_is_my_backup_alive(id, message_sendCh)
					}
					//time.Sleep(50 * time.Millisecond)
					//fmt.Println("Har jeg fått melding først?")
					if myBackupAlive == false {
						find_backup(id, p, &myBackupAlive, message_sendCh, &myBackupId)
					}
				}
				is_who_i_am_backup_for_gone(backupFor, p, n_to_p_task_manager)
				is_my_backup_gone(&myBackupAlive, p, &myBackupId)
				fmt.Println("Min id er:              ", id)
				fmt.Println("Min backupid er:        ", myBackupId)
				fmt.Println("Min myBackupAlive er:   ", myBackupAlive)
				fmt.Println("backupFor:              ", backupFor)
				//fmt.Println("")
			case <-time.After(time.Millisecond):
			}
			//time.Sleep(types.PAUSE_NET_LISTNER)
		}
	}()
	go func() {
		for {
			select {
			case a := <-message_receivedCh:
				if a.Destination == id || a.Destination == "broadcast" {
					switch a.Type & 224 {
					case types.ID_MODULE_DISTRIBUTOR:
						n_to_distri <- a

					case types.ID_MODULE_BACKUP_MANAGER:
						n_to_p_task_manager <- a

					case types.ID_MODULE_AMANAGER:
						n_to_a_tasks_manager <- a

					case types.ID_BACKUP_RESPONSE:
						n_to_a_tasks_manager2 <- a

					case types.ID_MODULE_NETWORK:
						//fmt.Println("NÅ kom meldingen:     ", a)
						message_receive_is_my_backup_alive(id, a, backupFor, message_sendCh)
						my_backup_is_alive(id, &myBackupAlive, a)
						backup_for(id, a, &backupFor)
						//fmt.Println("backupFor:              ", backupFor)
					}
				}
			case <-time.After(time.Millisecond):
			}
		}
	}()
}

//--------------------- Finner din backup  -----------------------------
func find_backup(id string, p peers.PeerUpdate, myBackupAlive *bool, message_sendCh chan types.MainData, myBackupId *string) {

	if len(p.Peers) > 1 { //xxxxxxxxxxxxxxxxxxxxxxxxxxx endre >= til > xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
		for {
			i := rand.Intn(len(p.Peers))
			if p.Peers[i] != id { //xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx endre == id til != id xxxxxxxxxxxxxxxxxxxxxxxxxx
				*myBackupAlive = true
				*myBackupId = p.Peers[i]

				message := types.MainData{}
				message.Source = id
				message.Destination = *myBackupId
				message.Type = types.YOU_ARE_MY_BACKUP
				row1 := []int{}
				row2 := []int{}
				message.Data = append(message.Data, row1)
				message.Data = append(message.Data, row2)
				message_sendCh <- message
				//fmt.Println("find_backup: ", message)
				//time.Sleep(50 * time.Millisecond)
				return
			}
			time.Sleep(types.PAUSE_NET_LISTNER)
		}
	}
}

//--------------------------- Legger hvem du er backupfor i en liste ----------------------
func backup_for(id string, a types.MainData, backupFor *[]string) {
	if (a.Destination == id) && ((a.Type & 31) == types.YOU_ARE_MY_BACKUP) {
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
		id = fmt.Sprintf("%s", localIP)
	}
	return id
}

//-----------------------  Melding for å sjekke om man har en backup --------------------
func send_message_is_my_backup_alive(id string, message_sendCh chan types.MainData) {
	message := types.MainData{}
	message.Source = id
	message.Destination = "broadcast"
	message.Type = types.IS_MY_BACKUP_ALIVE
	row1 := []int{}
	row2 := []int{}
	message.Data = append(message.Data, row1)
	message.Data = append(message.Data, row2)
	message_sendCh <- message
}

//--------------------------- Min Backup lever ------------------------------
func message_receive_is_my_backup_alive(id string, m types.MainData, backupFor []string, message_sendCh chan types.MainData) {
	if (m.Destination == "broadcast") && ((m.Type & 31) == types.IS_MY_BACKUP_ALIVE) {
		for i := range backupFor {
			if backupFor[i] == m.Source {
				message := types.MainData{}
				message.Source = id
				message.Destination = m.Source
				message.Type = types.IS_MY_BACKUP_ALIVE_TRUE
				row1 := []int{}
				row2 := []int{}
				message.Data = append(message.Data, row1)
				message.Data = append(message.Data, row2)
				message_sendCh <- message
				//fmt.Println("message_receive_backup_alive:    ", message)

			}
		}
	}
}

func my_backup_is_alive(id string, myBackupAlive *bool, m types.MainData) {
	if (m.Destination == id) && ((m.Type & 31) == types.IS_MY_BACKUP_ALIVE_TRUE) {
		*myBackupAlive = true
		//fmt.Println("my_backup_is_alive:    ", *myBackupAlive)
	}

}

func is_my_backup_gone(myBackupAlive *bool, p peers.PeerUpdate, myBackupId *string) {
	for j := range p.Lost {
		if *myBackupId == p.Lost[j] {
			*myBackupAlive = false
			*myBackupId = ""

		}
	}
}

func is_who_i_am_backup_for_gone(backupFor []string, p peers.PeerUpdate, n_to_p_task_manager chan types.MainData) {
	for i := range backupFor {
		for j := range p.Lost {
			if backupFor[i] == p.Lost[j] {
				message := types.MainData{}
				message.Source = p.Lost[j]
				message.Destination = ""
				message.Type = types.BACKUP_LOST
				row1 := []int{}
				row2 := []int{}
				message.Data = append(message.Data, row1)
				message.Data = append(message.Data, row2)
				n_to_p_task_manager <- message
			}
		}
	}
}
