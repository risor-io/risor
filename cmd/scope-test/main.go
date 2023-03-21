package main

import "fmt"

func main() {

	x := 42
	if x < 100 {
		x := 101
		fmt.Println("New X:", x)
	}

	fmt.Println("Done:", x)
}
