
package Network
 
//Message Constants
const (
		ID_MSG_TYPE_DISTIBUTOR_NEW_COMMAND = 24	//00011000
	)


type Commands struct { // Dette er dataen til elevator_handler som lagres i commands
    Next_floor int
    Current_floor int
    Person_inside bool

}

type NewOrder struct { // Dataen som kommer fra knappene
	Message_type int
    Floor int
    Direction int
}


type MainData struct { // Dette som sendes mellom heisene og legges i backoup
    Source string
    Destination string
    Message_type int
    Data [][]int
}
 