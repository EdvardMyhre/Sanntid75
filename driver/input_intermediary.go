package main  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import(
	"C"
	"fmt"
	"../Network"
)

//func Elev_get_button_signal(b Button){

//}

var b = Network.Button{}


func main(){
	b{
		Elev_id = 1
    	Button_type = 0
    	Floor 1	
	}
	
	
}

