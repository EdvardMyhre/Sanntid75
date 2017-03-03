package main
/*
#include<stdio.h>

hello()
{
    printf("Hello World!");
}
*/
import "C"
import(
	"fmt"
)

func main() {
	fmt.Println(C.hello())
	
}
