package fsm

import (
	"fmt"
	"time"

	"../elevController"
	"../elevio"
	"../lightController"
	"../orderHandler"
	"../statusManager"
)

//States for elevator
type State int

const (
	INIT        = 0
	IDLE        = 1
	EXECUTE     = 2
	LOST        = 3
	RESET       = 4
	OBSTRUCTION = 5
	STOP        = 6
)

func InitializeElev(addr string, ID int) {
	elevio.Init(addr, statusManager.GetNumFloors())
	statusManager.ResetAllOrders()
	elevio.SetMotorDirection(elevio.MD_Down)
	statusManager.MyElevInfo.Id = ID
	for elevio.GetFloor() != 0 {
	}
	elevio.SetFloorIndicator(0)
	elevio.SetMotorDirection(elevio.MD_Stop)
	//statusManager.SetMyElevInfo(0, currentOrder, state)
	fmt.Print("Elevator initialized with ID: ", ID, "\n")
}

func RunElevator() {
	// local runElev variables
	floor := 0
	state := IDLE
	NoOrder := statusManager.Order{Floor: -1, ButtonType: -1, Status: 0, Completed: false} // DEFINING noORDER
	currentOrder := NoOrder

	// CHANNELS FOR PASSING INFORMATION
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	motor_dir := make(chan int, 2)
	new_order := make(chan statusManager.Order, 100*statusManager.GetNumFloors()*statusManager.GetNumButtons())

	// GOROUTINES FOR GETTING INFORMATION
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	go elevController.HandleButtonPushed(drv_buttons, new_order)

	for {
		time.Sleep(20 * time.Millisecond)
		switch state {

		// STATE = IDLE
		case IDLE:
			fmt.Printf("IDLE\n")
			select {

			// IF A BUTTON IS PUSHED AND A NEW ORDER HAS ARRIEVED
			case currentOrder = <-new_order:
				currentOrder.Status = statusManager.GetMyElevInfo().Id // SETTING ORDER.STATUS=ELEV.ID
				statusManager.UpdateMyElevInfo(floor, currentOrder, state)

				// CHECKS IF ELEVATOR SHOULD TAKE ORDER
				if orderHandler.ShouldMyElevTakeOrder(currentOrder) {
					statusManager.UpdateOrder(int(currentOrder.ButtonType), currentOrder.Floor, statusManager.GetMyElevInfo().Id, true, false)
					motor_dir <- orderHandler.GetElevDirection(floor, currentOrder.Floor)
					state = EXECUTE
					statusManager.UpdateMyElevInfo(floor, currentOrder, state)
					break
				}
				statusManager.UpdateMyElevInfo(floor, NoOrder, state)

			//--TO HANDLE PUSHED STOP OR OBSTRUCTION--//
			case a := <-drv_obstr:
				if a {
					state = OBSTRUCTION
				}
			case a := <-drv_stop:
				if a {
					state = STOP
				}
			}
			//------------------------------------------//

		// STATE = EXECUTE
		case EXECUTE:
			fmt.Printf("EXECUTE\n")
			select {

			// WHEN ELEVATOR DRIVING
			case dir := <-motor_dir:
				elevController.SetMotorDirection(dir)
				if dir == 0 {
					lightController.SetLightDoorOpen()
					state = IDLE
					statusManager.UpdateMyElevInfo(floor, NoOrder, state)
				}

			// WHEN ELEVATOR AT FLOOR
			case floor = <-drv_floors:
				statusManager.UpdateMyElevInfo(floor, currentOrder, state)
				elevio.SetFloorIndicator(floor)
				if orderHandler.ShouldElevStopAtFloor(floor, currentOrder.Floor) {
					elevController.StopElevAtFloor(floor)
					dir := orderHandler.GetElevDirection(floor, currentOrder.Floor)
					if dir == 0 {
						state = IDLE
						statusManager.UpdateMyElevInfo(floor, NoOrder, state)
						for len(motor_dir) > 0 {
							<-motor_dir
						}
					} else {
						motor_dir <- dir
					}
				}

			//--TO HANDLE PUSHED STOP OR OBSTRUCTION--//
			case a := <-drv_obstr:
				if a {
					state = OBSTRUCTION
				}
			case a := <-drv_stop:
				if a {
					state = STOP
				}
			}
			//------------------------------------------//

		// STATE = OBSTRUCTION
		case OBSTRUCTION: // Nå er obstruction true
			fmt.Printf("OBSTRUCTION\n")
			elevio.SetMotorDirection(elevio.MD_Stop)
			if statusManager.CheckForOrdersAtCurrentFloor() {
				elevController.StopElevAtFloor(floor)
			}

			select {
			case obs := <-drv_obstr: // Dette er tilfelle obstruciton får inn et nytt signal
				if !obs { // Hvis obs is flase, send til EXECUTE or IDLE
					if statusManager.CheckForOrders() {
						fmt.Printf("checking for orders \n")
						motor_dir <- orderHandler.GetElevDirection(floor, currentOrder.Floor)
						state = EXECUTE
					} else {
						state = IDLE
					}
				}
			}

		// STATE = STOP
		// Funker ikke helt som den skal, du må holde inne STOP (p) for at den skal slette alle ordre,
		// funker ikke bare å trykke på STOP (p)
		case STOP:
			fmt.Printf("STOP\n")
			elevio.SetMotorDirection(elevio.MD_Stop)
			statusManager.ResetAllOrders()
			elevio.SetStopLamp(true)
			select {
			case stp := <-drv_stop:
				if stp {
					elevio.SetMotorDirection(elevio.MD_Stop)

				} else {
					elevio.SetStopLamp(false)
					state = IDLE
				}
			}
		}
	}
}
