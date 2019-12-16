package main

import (
"fmt"
"time"
"math/rand"
"ai"
)

/*
Get 
Steps,id.  
return 
1. situation judgement、 lose or win
2. get next step( with id)
*/

type StepsInfo struct{
	x,y []int
	forbid bool
	level int
}

func CreateRand() int64{
	time.Seed(time.Now().UnixNano())
	return rand.Int63()
}

func ReplyOver(over int){
}

// post setsteps, 
// if over return result/reset
// if not over (go getstep, store id) return situation and id
// post id get step and situation(and result)

func CreateBySteps(info *StepsInfo)(*AIPlayer,id int64){
	cur:=len(info)
	p,_:=ai.InitPlayer(ai.WHITE,info.level,info.forbid)
	id=CreateRand()
	if cur%2 ==0{// next is black, so ai use black
		p.human=WHITE
		p.robot=BLACK
	}
	for k,v:=range info.x{
		p.SetStep(v,p.y[k])
		if over:=p.IsOver();over!=0{
			ReplyOver(over)
			break
		}
	}
	return p,id
}

func main(){
/*    fmt.Println("Start:")
    player,err:=ai.InitPlayer(ai.WHITE,0)
	if err!=nil{
        fmt.Println("Init server error:",err)
    }
	player.SetStep(7,7)*/
	
}

