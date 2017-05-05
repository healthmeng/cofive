package main

import (
"fmt"
"ai"
)

func draw(player *ai.AIPlayer){
	frame:=player.GetFrame()
	fmt.Print("  ")
	for i:=0;i<15;i++{
		fmt.Printf("%2d",i)
	}
	fmt.Println("")
	for i:=0;i<15;i++{
		fmt.Printf("%-2d",i)
		for j:=0;j<15;j++{
			switch frame[j][i]{
			case 0:
				fmt.Printf(" .")
			case 1:
				fmt.Printf(" o")
			case 2:
				fmt.Printf(" x")
			}
		}
		fmt.Println("")
	}
}

func main(){

    fmt.Println("Start:")
    player,err:=ai.InitPlayer(ai.WHITE,0)
	if err!=nil{
        fmt.Println("Init server error:",err)
		return
    }
	draw(player)
	over:=0
	for ;over==0;over=player.IsOver(){
		var x,y int
		fmt.Scanf("%d%d",&x,&y)
		player.SetStep(x,y)
		draw(player)
		if over=player.IsOver();over!=0{
			break
		}
		x,y=player.GetStep()
		draw(player)
	}
	if over==1{
		fmt.Println("Black win")
	}else{
		fmt.Println("White win")
	}
}

