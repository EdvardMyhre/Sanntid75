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



//-----------------------  Finner lokal ip og deklarerer myBackupId --------------------
	id := find_localip()
	myBackupId := ""
	myBackupAlive := false
	var backupFor []string
	backupFor = ""


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

			if message_send.Destination == "backup"{
				if myBackupAlive == false{
					myBackupId = find_backup()
				}
				message_send.Destination = myBackupId
			}
			message_sendCh <- message_send
		}
	}()






	//-----------------------  Fordeler det som kommer fra broadcast og håndterer hvem som er i livet --------------------
	fmt.Println("Started")
	go func(){
		for {
			select {
			case p := <-peerUpdateCh:
				if myBackupId == ""{
					for i := 0; i < 5; i++{
						send_message_is_my_backup_alive(id, message_sendCh)
					}
					time.Sleep(200 * time.Millisecond)
					if myBackupAlive == false {
						myBackupId = find_backup(id, p, message_sendCh)
					}
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
						message_receive_backup_alive(id, a, backupFor, message_sendCh)
						my_backup_is_alive(id, &myBackupAlive)
						backup_for(a, &backupFor)
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



//--------------------- Finner din backup  -----------------------------
func find_backup(id string, p peers.PeerUpdate, message_sendCh chan structer.MainData) string {
	if len(p.Peers) > 1 {
		for {
			i := rand.Intn(len(p.Peers))
			//fmt.Println(i)
			myBackupId := p.Peers[i]
			if myBackupId != id{
				myBackupAlive = true

				message := structer.MainData{}
				message.Source = id
				message.Destination = m.myBackupId
				message.Message_type = messageid.ID_MSG_TYPE_YOU_ARE_MY_BACKUP
				row1 := []int{}
				row2 := []int{}
				message.Data = append(message.Data, row1)
				message.Data = append(message.Data, row2)
				message_sendCh <- message

				return myBackupId
			}
		}
	}
	return ""
}






//--------------------------- Legger hvem du er backupfor i en liste ----------------------
func backup_for(id string,a structer.MainData, *backupFor) {
	if (a.Destination == id) && ((a.Message_type & 31) == messageid.ID_MSG_TYPE_YOU_ARE_MY_BACKUP){
		*backupFor = append(*backupFor, a.Source)
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
	/*for {
		message_sendCh <- message
		time.Sleep(1 * time.Second)
	}*/
}()




//--------------------------- Min Backup lever ------------------------------
func message_receive_backup_alive(id string, m structer.MainData, backupFor []string, message_sendCh chan structer.MainData){
	if (m.Destination == "broadcast") && ((m.Message_type & 31) == messageid.ID_MSG_TYPE_IS_MY_BACKUP_ALIVE){
		for i := range backupFor{
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

			}
		}
	}

}

func my_backup_is_alive(id string,*myBackupAlive bool) {
	if (m.Destination == id) && ((m.Message_type & 31) == messageid.ID_MSG_TYPE_IS_MY_BACKUP_ALIVE_TRUE){
		*myBackupAlive = true
	}

}




