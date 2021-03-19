package main

import (
    "fmt"
    "time"

    "./elevio"
    "./fsm"
    //"./orderHandler"
    
)


func main(){
    const numFloors = 4
    const numButtons =3

   
    

    
    var d elevio.MotorDirection = elevio.MD_Stop
    var currentFloor int = 0
    var state fsm.State

    // Initialize elevator
    fsm.InitializeElev("localhost:12345", numFloors, numButtons)
    
    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)    
    
    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)

   
    
    
    for {
        time.Sleep(20 * time.Millisecond)
        switch fsm.MyElevInfo.State{
        case fsm.IDLE:
            fmt.Printf("IDLE\n")
        select {
            case a := <- drv_buttons:
                fmt.Printf("%+v\n", a)
                fsm.SetElevLights(a.Button, a.Floor, true)
                fsm.UpdateOrder(a.Button, a.Floor, 1)
                if fsm.CheckForOrdersAbove(fsm.MyElevInfo.Floor){
                    d=elevio.MD_Up
                } else if fsm.CheckForOrdersBelow(fsm.MyElevInfo.Floor){
                    d=elevio.MD_Down
                }
                if d != elevio.MD_Stop{
                    elevio.SetMotorDirection(d)
                    fsm.SetMyElevInfo(currentFloor, d, fsm.EXECUTE)
                }
            }
        case fsm.EXECUTE:
                fmt.Printf("EXECUTE\n")
            select{
            case b := <- drv_buttons:
                fmt.Printf("%+v\n", b)
                fsm.SetElevLights(b.Button, b.Floor, true)
                fsm.UpdateOrder(b.Button, b.Floor, 1)

            case a := <- drv_floors:
                fmt.Printf("fFIIIITTTE\n")
               
                elevio.SetFloorIndicator(a)
                
                fmt.Printf("%+v\n", a)
            
                    // Looking for orders in the same direction
                d = fsm.GetMdirElev()
                fmt.Printf("%+v\n", elevio.MotorDirection(d))
                if d==elevio.MD_Down{
                    if !fsm.CheckForOrdersBelow(a){
                        fmt.Printf("STOP0")
                        d=elevio.MD_Stop
                        //fsm.SetMyElevInfo(a, fsm.IDLE)
                        state=fsm.IDLE
                    }
                } else if d==1{
                    if !fsm.CheckForOrdersAbove(a){
                        fmt.Printf("STOP1")
                        d=elevio.MD_Stop
                        //fsm.SetMyElevInfo(a, fsm.IDLE)
                        state=fsm.IDLE
                    }
                } 
                    if fsm.CheckForOrdersAbove(a){
                        d=elevio.MD_Up
                        state = fsm.EXECUTE
                        fmt.Printf("STOP2")
                    }else if  fsm.CheckForOrdersBelow(a){
                        d=elevio.MD_Down
                        state=fsm.EXECUTE
                        fmt.Printf("STOP3")
                
                

            }
            elevio.SetFloorIndicator(a)
            fmt.Printf("STOP4")
        
                    /*
                    if a == numFloors-1 {
                        d = elevio.MD_Down
                    } else if a == 0 {
                        d = elevio.MD_Up
                    }
                    */
                    var checkIfThereAreOrders int = 0
                    for b:=0; b<numButtons; b++{
                        
                        if (fsm.MyElevInfo.Orders[a][b].Status==1){
                            checkIfThereAreOrders +=1
                            elevio.SetMotorDirection(elevio.MD_Stop)
                            
                            fsm.UpdateOrder(elevio.ButtonType(b), a, 0)
                            fsm.SetElevLights(elevio.ButtonType(b), a, false)
                            fmt.Print("TakingOrder\n")
                            
                        }
                    }
                    if checkIfThereAreOrders !=0{
                        elevio.SetDoorOpenLamp(true)
                        time.Sleep(3*time.Second)
                        elevio.SetDoorOpenLamp(false)
                    }
                    
                    if !fsm.CheckForOrders(){
                        state=fsm.IDLE
                        d = elevio.MD_Stop
                    }
                    elevio.SetMotorDirection(d)
                    fsm.SetMyElevInfo(a, d, state)
            
        
          
            

            
    
            
            
        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }

        // STOP BUTTON
        case a := <- drv_stop:
            elevio.SetStopLamp(true)
            fmt.Printf("%+v\n", a)
            fmt.Printf("stoppppping\n")
            for f := 0; f < numFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
                
            } else {
                elevio.SetStopLamp(false)
                elevio.SetMotorDirection(d)
            }
           
            
        }
        
    }
    }    
}


