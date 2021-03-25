package statusManager

import (
	"fmt"

	"../lightController"
)

const numFloors = 4
const numButtons = 3

type Order struct {
	Floor      int
	ButtonType int
	Status     int
	Confirmed  bool
	Completed  bool
}

type OrderStatus int

const (
	TimedOut OrderStatus = -1 //Gått for lang tid at den ikke er tatt
	Inactive             = 0  //btn not pushed for this order
	Pending              = 1  //There is an order here, waiting
	Active               = 2  // An elev is on their way
)

type Elev struct {
	Id           int
	Floor        int
	CurrentOrder Order
	State        int
	Orders       [numFloors][numButtons]Order
}

// INFORMATION ABOUT MYSELF
var MyElevInfo Elev

// INFORMATION ABOUT THE OTHER ELEVATORS ON THE NETWORK
var OtherElevInfo []Elev // not used yet

func ResetAllOrders() {
	for f := 0; f < numFloors; f++ {
		for b := 0; b < numButtons; b++ {
			MyElevInfo.Orders[f][b].Floor = f
			MyElevInfo.Orders[f][b].ButtonType = b
			MyElevInfo.Orders[f][b].Status = 0
			MyElevInfo.Orders[f][b].Confirmed = false
			MyElevInfo.Orders[f][b].Completed = false
		}
	}
	lightController.InitializeLights(numFloors, numButtons)
}

func UpdateMyElevInfo(floor int, currentOrder Order, state int) {
	MyElevInfo.Floor = floor
	MyElevInfo.State = state
	MyElevInfo.CurrentOrder = currentOrder
}

func UpdateOrder(button int, floor int, status int, confirmed bool, completed bool) {
	fmt.Printf("UPDATE ORDER %d\n", floor)
	MyElevInfo.Orders[floor][button].Status = status
	MyElevInfo.Orders[floor][button].Confirmed = confirmed
	MyElevInfo.Orders[floor][button].Completed = completed

	// SETTING LIGTHS (litt jalla løsning, men kanskje dette kan gjøres om)
	var set bool = false
	if status != 0 {
		set = true
	}
	lightController.SetOrderLights(button, floor, set)
}

//-------------------------------------------------------------//
//----------------------GETTING FUNCTIONS----------------------//
//-------------------------------------------------------------//
func GetMyElevInfo() Elev {
	return MyElevInfo
}

func GetOrder(floor int, btn int) Order {
	return MyElevInfo.Orders[floor][btn]
}

func GetNumFloors() int {
	return numFloors
}

func GetNumButtons() int {
	return numButtons
}

//--------------------------------------------------------------//
//----------------------CHECKING FUNCTIONS----------------------//
//--------------------------------------------------------------//
// USED FOR SINGLE ELEVATOR, NOT MULTIPLE
func CheckForOrdersBelow(myCurrentFloor int) bool {
	for f := 0; f < myCurrentFloor; f++ {
		for b := 0; b < numButtons; b++ {

			if MyElevInfo.Orders[f][b].Status == 1 {
				return true
			}
		}
	}
	return false
}

// USED FOR SINGLE ELEVATOR, NOT MULTIPLE
func CheckForOrdersAbove(myCurrentFloor int) bool {
	if myCurrentFloor == 3 {
		return false
	}

	for f := myCurrentFloor + 1; f < numFloors; f++ {

		for b := 0; b < numButtons; b++ {

			if MyElevInfo.Orders[f][b].Status == 1 {
				return true
			}
		}
	}
	return false
}

// CHECK IF THERE ARE ANY ORTHERS IN MyElevInfo
func CheckForOrders() bool {
	for f := 0; f < numFloors; f++ {
		for b := 0; b < numButtons; b++ {

			if MyElevInfo.Orders[f][b].Status == 1 {
				return true
			}
		}
	}
	return false
}

// CHECKS IF THERE ARE ANY ORDERS AT MyElevInfo.floor (current floor)
func CheckForOrdersAtCurrentFloor() bool {
	for b := 0; b < numButtons; b++ {
		if MyElevInfo.Orders[MyElevInfo.Floor][b].Status != 0 {
			return true
		}
	}
	return false
}
