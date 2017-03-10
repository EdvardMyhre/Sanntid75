package types

import "time"

type Task struct{
	Type int
	Floor int
	Add int
}

type Button struct {
    Type int
    Floor int
    Instant time.Time		
}

type Status struct {
	Destination_floor int
	Floor int
	Prev_floor int
	Finished int
}

type MainData struct { // Dette som sendes mellom heisene og legges i backoup
    Source string
    Destination string
    Message_type int
    Data [][]int
}