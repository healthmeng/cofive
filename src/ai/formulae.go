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
	return nil
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
