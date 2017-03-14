package types

import "time"

type Task struct {
	Type     int
	Floor    int
	Finished int
	Assigned int
}

type Button struct {
	Type    int
	Floor   int
	Instant time.Time
}

type Status struct {
	Destination_floor int
	Floor             int
	Prev_floor        int
	Finished          int
	Between_floors    int
}

type MainData struct {
	Source      string
	Destination string
	Type        int
	Data        [][]int
}
