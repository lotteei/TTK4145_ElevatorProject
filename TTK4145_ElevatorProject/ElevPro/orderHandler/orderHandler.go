package orderHandler

import (
	"fmt"

	"../statusManager"
)

// GETS ELEVATOR DIRECTION FROM WHERE IT WANTS TO GO
func GetElevDirection(currentFloor int, wantedFloor int) int {
	fmt.Printf("GET ELEV DIR: ")
	if wantedFloor == -1 || currentFloor == wantedFloor {
		fmt.Printf("0\n")
		return 0
	} else if wantedFloor < currentFloor {
		fmt.Printf("-1\n")
		return -1
	} else {
		fmt.Printf("1\n")
		return 1
	}
}

// IKKE HELT FERDIG. DENNE RETURNERER KUN true NÅR DET ER EN CAB-ORDER
// LEGG TIL EN KOSTFUNKSJON ELLER NOE TIL Å VURDERE HALL-ORDERS
func ShouldMyElevTakeOrder(order statusManager.Order) bool {
	if order.ButtonType == 2 {
		return true
	}

	// Ta INN EN COSTFUNKSJON ELLER NOE
	return false
}

// CHECK IF ELEVATOR SHOULD STOP AT A FLOOR
func ShouldElevStopAtFloor(currentFloor int, wantedFloor int) bool {
	dir := GetElevDirection(currentFloor, wantedFloor)
	if dir == 0 {
		return true
	}
	if statusManager.GetOrder(currentFloor, 0).Status == 1 && dir == 1 {
		return true
	}
	if statusManager.GetOrder(currentFloor, 1).Status == 1 && dir == -1 {
		return true
	}
	if statusManager.GetOrder(currentFloor, 2).Status == 1 {
		return true
	}
	return false

}
