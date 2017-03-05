package types

import "time"

type Button struct { // Data from buttons are stored here
    Type int
    Floor int
    Instant time.Time		
}

type Status struct {
	Floor int
	Prev_floor int
	Finished int
}