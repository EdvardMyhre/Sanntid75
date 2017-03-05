package types

type Task struct { // Data from buttons are stored here
    Elev_id int
    Button_type int
    Floor int		
}

type Status struct {
	Floor int
	Prev_floor
	Finished int
}