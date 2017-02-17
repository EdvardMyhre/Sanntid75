package main

import (
    "fmt"
)

type Command struct{
    Elevator int
    Next_floor int
}


/*func (p Command) String() string {
   return fmt.Sprintf("%v %v Next_floor", p.Elevator, p.Next_floor)
}*/


func main() {
    a := Command{1,3};
    fmt.Println("Heisen som skal kjøre er nr: ",a.Elevator,", og den skal til ",a.Next_floor," etasje");
}
