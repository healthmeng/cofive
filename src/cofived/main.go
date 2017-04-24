package main

import (
"fmt"
"ai"
)
func main(){

    fmt.Println("Start:")
    player,err:=ai.InitPlayer(ai.WHITE,0)
	if err!=nil{
        fmt.Println("Init server error:",err)
    }
	player.SetStep(7,7)
}

