package structer

type Command struct { // Dette er dataen til elevator_handler som lagres i commands
	Next_floor    int
	Current_floor int
	Person_inside bool
}

type InternalMessage struct { // Mellom moduler (ikke network)
	Message_type int
	Data         []int
}

type Button struct { // Data from buttons are stored here
	Elev_id     int
	Button_type int
	Floor       int
}

type MainData struct { // Dette som sendes mellom heisene og legges i backoup
	Source       string
	Destination  string
	Message_type int
	Data         [][]int
}

//type NewOrder struct { // Dataen som kommer fra knappene
//    Floor int
//    Direction int
//}
