package elevator

import "../driver"
import "../types"
import "time"



func Controller(statusc chan<- types.Status, taskc <-chan int){
	driver.SetMotorDirection(types.MOTOR_DIR_DOWN)
	for driver.GetFloorSensorSignal() != 0{
	}
	driver.SetMotorDirection(types.MOTOR_DIR_STOP)

	status := types.Status{Floor: 0, Prev_floor: 1, Finished: 1}
	destination_floor := 0
	floor_signal := 0
	driver.SetFloorIndicator(0)

	for{
		floor_signal = driver.GetFloorSensorSignal()
		if floor_signal >= 0 && floor_signal != status.Floor{
			status.Prev_floor = status.Floor
			status.Floor = floor_signal
			driver.SetFloorIndicator(status.Floor)
			statusc <- status 								//block trøbbel?
		}

		if status.Finished == 1{
			destination_floor = <- taskc
			status.Finished = 0
			statusc <- status
		}

		if status.Finished == 0{
			if floor_signal == destination_floor{
				driver.SetMotorDirection(types.MOTOR_DIR_STOP)
				driver.SetDoorOpenLamp(types.LAMP_ON)
				time.Sleep(time.Second * 4)
				driver.SetDoorOpenLamp(types.LAMP_OFF)
				status.Finished = 1
				statusc <- status //block trøbbel?
			} else if floor_signal >= 0 && floor_signal < destination_floor{
				driver.SetMotorDirection(types.MOTOR_DIR_UP)
			} else if floor_signal >= 0 && floor_signal > destination_floor{
				driver.SetMotorDirection(types.MOTOR_DIR_DOWN)
			}
		}
	}
	
}


func ButtonPoller(buttonc chan<- types.Button){
	//Initialising variables
	button := types.Button{}
	var button_pushes []types.Button
	var button_pushes_this_loop []types.Button

	button_type := 0

	button.Instant = time.Now()
	for button_type = 0; button_type < 3; button_type++{
		for floor := 0; floor < types.NUMBER_OF_FLOORS; floor++{
			button.Type = button_type
			button.Floor = floor
			button_pushes = append(button_pushes, button)
		}
	}

	//Polling buttons
	for{
		button_pushes_this_loop = nil

		button_type = types.BTN_TYPE_UP
		for floor := 0; floor < types.NUMBER_OF_FLOORS-1; floor++{
			if driver.GetButtonSignal(button_type, floor) != 0{
				button.Type = button_type
				button.Floor = floor
				button.Instant = time.Now()
				button_pushes_this_loop = append(button_pushes_this_loop, button)
			}
		}
		button_type = types.BTN_TYPE_DOWN
		for floor := 1; floor < types.NUMBER_OF_FLOORS; floor++{
			if driver.GetButtonSignal(button_type, floor) != 0{
				button.Type = button_type
				button.Floor = floor
				button.Instant = time.Now()
				button_pushes_this_loop = append(button_pushes_this_loop, button)
			}
		}
		button_type = types.BTN_TYPE_COMMAND
		for floor := 0; floor < types.NUMBER_OF_FLOORS; floor++{
			if driver.GetButtonSignal(button_type, floor) != 0{
				button.Type = button_type
				button.Floor = floor
				button.Instant = time.Now()
				button_pushes_this_loop = append(button_pushes_this_loop, button)
			}
		}
		

		for i := 0; i < len(button_pushes_this_loop); i++{
			for j := 0; j < len(button_pushes); j++{
				if button_pushes_this_loop[i].Type == button_pushes[j].Type && button_pushes_this_loop[i].Floor == button_pushes[j].Floor{
					if time.Since(button_pushes[j].Instant).Seconds() - time.Since(button_pushes_this_loop[i].Instant).Seconds() > 2{
						button_pushes[j] = button_pushes_this_loop[i]
						buttonc <- button_pushes_this_loop[i]
					}
				}
			}
		}
	}
}


















//TEST CONTROLLER
// func main() {

// 	statusc := make(chan types.Status)
// 	taskc := make(chan int)
// 	go elevator.Controller(statusc, taskc)

// 	for{
// 		select{
// 			case status := <- statusc:
// 				fmt.Println("Floor : ", status.Floor, " Prev_floor: ", status.Prev_floor, " Finished: ", status.Finished)
// 			case <- time.After(time.Second):
// 		}
// 		if rand.Intn(10) == 0{
// 			destination_floor := rand.Intn(4)
// 			select{
// 			case taskc <- destination_floor:
// 				fmt.Println("Destination floor: ", destination_floor)
// 			default:
// 			}
// 		}

// 	}
// }




//TEST BUTTONS
// func main() {
// 	driver.Init()
// 	buttonc := make(chan types.Button)
// 	go elevator.Buttons(buttonc)
// 	for{
// 		button := <- buttonc
// 		fmt.Println("Type: ", button.Type, " Floor: ", button.Floor)
// 	}

// }