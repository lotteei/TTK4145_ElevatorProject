package elevController

import (
	"fmt"
	"time"

	"../elevio"
	"../lightController"
	"../statusManager"
)

func SetMotorDirection(MotorDir int) {
	elevio.SetMotorDirection(elevio.MotorDirection(MotorDir))
}

// HANDLES ALL THE ACTIONS THAT MUST BE DONE WHEN A BUTTON IS PUSHED
func HandleButtonPushed(btn_pressed chan elevio.ButtonEvent, new_order chan<- statusManager.Order) {
	fmt.Printf("HANDLE BUTTON PUSHED\n")
	for {
		time.Sleep(20 * time.Millisecond)
		select {
		case btn := <-btn_pressed:
			fmt.Printf("%+v\n", btn)
			order := statusManager.GetOrder(btn.Floor, int(btn.Button)) //GETS THE ORDER FROM ORDERLIST WHERE THE BUTTON IS PUSHED
			if order.Status == 0 {                                      // IF ORDER IS INACTIVE (0)
				statusManager.UpdateOrder(order.ButtonType, order.Floor, 1, false, false) // UPDATE ORDER TO PENDING
				if order.ButtonType == 2 {                                                // IF ORDER IS CAB ORDER (2)
					new_order <- order // SETS THE ORDER TO new_order CHANNEL
				}
			}
		}
	}
}

func StopElevAtFloor(floor int) {
	fmt.Printf("STOPPIG AT FLOOR \n")
	for b := 0; b < statusManager.GetNumButtons(); b++ {
		// CHECKS  STATUS ON ORDERS ON FLOOR
		status := int(statusManager.GetOrder(floor, b).Status)
		// IF STATUS=PENDING OR STATUS=ELEV.ID, SET ORDER AS COMPLETED
		if status == 1 || status == statusManager.GetMyElevInfo().Id {
			statusManager.UpdateOrder(b, floor, status, false, true)
			lightController.SetOrderLights(b, floor, false)
		}
	}

	// SET DOORLIGHT
	lightController.SetLightDoorOpen()

	// INITIALIZE ORDERS, THIS WAY THEY CAN BE TAKEN AGAIN
	for b := 0; b < statusManager.GetNumButtons(); b++ {
		status := int(statusManager.GetOrder(floor, b).Status)
		if status == 1 || status == statusManager.GetMyElevInfo().Id || statusManager.GetOrder(floor, b).Completed == true {
			statusManager.UpdateOrder(b, floor, 0, false, false)
		}
	}
}
