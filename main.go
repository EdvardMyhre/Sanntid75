package main

import (
	"./NetworkModul"
	"./NetworkModul/network/structer"
	"fmt"
)

func main() {
	network_to_distributing_state_machine := make(chan structer.MainData)
	go network.Network_start(network_to_distributing_state_machine)
	for {

		select {
		case p := <-network_to_distributing_state_machine:
			fmt.Println("Sendt fra main:  ", p)
		}
	}

}
