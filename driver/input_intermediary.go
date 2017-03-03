package main
/*
#include<stdio.h>

hello()
{
    return 1;
}
*/
import "C"
import(
	"fmt"
)

func main() {
	fmt.Println(C.hello())
	
}
