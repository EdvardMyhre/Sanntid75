package input

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "../driver/elev.h"
*/
import "C"

func Cab_4() int {
	a := C.int(1)
	a = C.elev_get_button_signal(C.int(2), C.int(4))
	return int(a)

}
