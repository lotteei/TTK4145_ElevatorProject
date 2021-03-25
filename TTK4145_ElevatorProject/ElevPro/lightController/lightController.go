package lightController

import (
	"time"

	"../elevio"
)

func InitializeLights(numFloors int, numButtons int) {
	for f := 0; f < numFloors; f++ {
		for b := 0; b < numButtons; b++ {
			SetOrderLights(b, f, false)
		}
	}
}

// SETS ORDER LIGHTS
func SetOrderLights(button int, floor int, set bool) {
	elevio.SetButtonLamp(elevio.ButtonType(button), floor, set)
}

// SETS DOOR OPEN FOR 3 SECONDS
func SetLightDoorOpen() {
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)
	time.Sleep(3 * time.Second)
	elevio.SetDoorOpenLamp(false)
}
