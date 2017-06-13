package ai

import (
_"log"
"math/rand"
"time"
)

func (player* AIPlayer)getWhiteFormula(x,y int)*StepInfo{
	sts:=make([]StepInfo,0,8)
	for i:=-1;i<=1;i++{
		for j:=-1;j<=1;j++{
			if player.frame[x+i][y+j]==0{
				sts=append(sts,StepInfo{x+i,y+j,2})
			}
		}
	}
	nst:=len(sts)
	rnd:=rand.New(rand.NewSource(time.Now().UnixNano()))
	return &sts[rnd.Int()%nst]
}

func (player* AIPlayer)getThirdStep() *StepInfo{
    sts := make([]StepInfo, 0, MAX_STEP)
	x:= player.steps[0].x
	y:= player.steps[0].y
	var orders [5]int=[5]int{-1,1,0,-2,2}
	for i:= 0 ;i<5;i++  {
		tmpx:=orders[i]+x
	    if tmpx< 0 || tmpx >= 15 {
	        continue
	    }
	    for j := 0;j<5; j++ {
			tmpy:=orders[j]+y
	        if tmpy < 0 || tmpy >= 15 {
	            continue
	        }
	        if player.frame[tmpx][tmpy] == 0 {
				sts = append(sts, StepInfo{tmpx, tmpy, 1})
	        }
	    }
	}
	nst:=len(sts)
	rnd:=rand.New(rand.NewSource(time.Now().UnixNano()))
	return &sts[rnd.Int()%nst]
}

func (player* AIPlayer)IsFormula() bool{
	dx:=player.steps[0].x-player.steps[1].x
	dy:=player.steps[0].y-player.steps[1].y
	if dx<=1 && dx>=-1 && dy<=1 && dy>=-1{
		return true
	}
	return false
}

func (player* AIPlayer)TryFormula() *StepInfo{
	nst:=player.curstep
	switch nst{
	case 0:
		return &StepInfo{x:7,y:7,bw:1}
	case 1:
		return player.getWhiteFormula(player.steps[0].x,player.steps[0].y)
	case 2:
		if player.IsFormula(){
			return player.getThirdStep()
		}
	}
	return nil
}
