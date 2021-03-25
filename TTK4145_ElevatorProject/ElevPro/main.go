package main

import (
	"./fsm"
	//"./orderHandler"
)

func main() {

	// Initialize elevator
	fsm.InitializeElev("localhost:12345", 1)

	// Run elevator
	go fsm.RunElevator()
	select {}
}
