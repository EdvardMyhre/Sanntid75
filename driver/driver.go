package driver

// #include "elev.h"
// #cgo CFLAGS: -std=gnu11
// #cgo LDFLAGS: -L. -lcomedi -lm
import "C"

import (
	"fmt"

	"../config"
	"../types"
)

func Init() {
	switch config.ElevatorType {
	case types.ET_SIM:
		C.elev_init(C.ET_Simulation)
		break
	case types.ET_COMEDI:
		C.elev_init(C.ET_Comedi)
		break
	}
}

func SetMotorDirection(direction int) {

	switch direction {
	case types.MOTOR_DIR_UP:
		C.elev_set_motor_direction(C.DIRN_UP)
		break

	case types.MOTOR_DIR_DOWN:
		C.elev_set_motor_direction(C.DIRN_DOWN)
		break

	case types.MOTOR_DIR_STOP:
		C.elev_set_motor_direction(C.DIRN_STOP)
		break

	}

}

func SetButtonLamp(buttonType int, floor int, value C.int) {

	// Cast Go ints to C ints
	c_floor := C.int(floor)
	c_value := C.int(value)

	switch buttonType {
	case types.BTN_TYPE_UP:
		C.elev_set_button_lamp(C.BUTTON_CALL_UP, c_floor, c_value)
		break

	case types.BTN_TYPE_DOWN:
		C.elev_set_button_lamp(C.BUTTON_CALL_DOWN, c_floor, c_value)
		break

	case types.BTN_TYPE_COMMAND:
		C.elev_set_button_lamp(C.BUTTON_COMMAND, c_floor, c_value)
		break

	}

}

func SetFloorIndicator(floor int) {
	c_floor := C.int(floor)
	C.elev_set_floor_indicator(c_floor)
}

func SetDoorOpenLamp(value C.int) {
	C.elev_set_door_open_lamp(value)
}

func GetButtonSignal(buttonType, floor int) int {

	c_floor := C.int(floor)

	switch buttonType {
	case types.BTN_TYPE_UP:
		return int(C.elev_get_button_signal(C.BUTTON_CALL_UP, c_floor))

	case types.BTN_TYPE_DOWN:
		return int(C.elev_get_button_signal(C.BUTTON_CALL_DOWN, c_floor))

	case types.BTN_TYPE_COMMAND:
		return int(C.elev_get_button_signal(C.BUTTON_COMMAND, c_floor))

	default:
		fmt.Println("invalid button")
		return 0
	}

}

func GetFloorSensorSignal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func GetStopSignal() int {
	return int(C.elev_get_stop_signal())
}
