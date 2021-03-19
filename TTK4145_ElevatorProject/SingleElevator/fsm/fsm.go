package fsm

import (
	"fmt"
	"../elevio"
)

const numFloors = 4
const numButtons =3


//States for elevator 
type State int
const (
	INIT  State  = 0
	IDLE    = 1
	EXECUTE = 2
	LOST    = 3
	RESET   = 4
)

type Order struct {
	Floor		int
	ButtonType 	int
	Status      int //1 finished, 0 not finished 
}
type OrderStatus int

const(
	Pending OrderStatus = 1
	Inactive            = 0
)

type Elev struct {
	Id           int
	Floor        int
	CurrentOrder Order
	State        State
	Mdir 	 	 elevio.MotorDirection
	Orders       [numFloors][numButtons]Order
}

var MyElevInfo Elev



func InitializeElev(addr string, num_Floors int, num_Buttons int){
	elevio.Init(addr, num_Floors)
	for f := 0; f<numFloors; f++ {
		for b:=0; b<numButtons; b++{
			MyElevInfo.Orders[f][b].Status=0
			elevio.SetButtonLamp(elevio.ButtonType(b), f, false)
		}
	}
	elevio.SetMotorDirection(elevio.MD_Down)
	MyElevInfo.Floor=0
    for elevio.GetFloor() != 0{ 
    }
	elevio.SetFloorIndicator(0)
    elevio.SetMotorDirection(elevio.MD_Stop)
	MyElevInfo.State=IDLE
	fmt.Print("Elevetor initialized\n")
}

func UpdateOrder(button elevio.ButtonType, floor int, status int){
	MyElevInfo.Orders[floor][button].Status = status
}

func SetElevLights(button elevio.ButtonType, floor int, set bool){
	elevio.SetButtonLamp(button, floor, set)
}
func GetMdirElev() elevio.MotorDirection{
	return MyElevInfo.Mdir
}

func SetMyElevInfo( floor int, dir elevio.MotorDirection, state State){
	MyElevInfo.Floor=floor
	MyElevInfo.Mdir=dir
	MyElevInfo.State=state
}
//func addOrder(ButtonPress chan elevio.ButtonEvent, lightsChannel chan<- elevio.PanelLight){
	
func CheckForOrdersBelow(myCurrentFloor int) bool{
	for f := 0; f<myCurrentFloor; f++ {
		for b:=0; b<numButtons; b++{

			if MyElevInfo.Orders[f][b].Status==1{
				return true
			}
		}
	}
	return false
}
func CheckForOrdersAbove(myCurrentFloor int) bool{
	for f := myCurrentFloor; f<numFloors; f++ {
		for b:=0; b<numButtons; b++{

			if MyElevInfo.Orders[f][b].Status==1{
				return true
			}
		}
	}
	return false
}
func CheckForOrders() bool{
	for f := 0; f<numFloors; f++ {
		for b:=0; b<numButtons; b++{

			if MyElevInfo.Orders[f][b].Status==1{
				return true
			}
		}
	}
	return false
}