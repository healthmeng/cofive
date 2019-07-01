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
			if x+i>=0 && x+i<15 && y+j>=0 && y+j<15 && player.frame[x+i][y+j]==0{
				sts=append(sts,StepInfo{x+i,y+j,2,false})
			}
		}
	}
	nst:=len(sts)
	rnd:=rand.New(rand.NewSource(time.Now().UnixNano()))
	return &sts[rnd.Int()%nst]
}

func (player* AIPlayer)getThirdStep() *StepInfo{
    sts := make([]StepInfo, 0, MAX_STEP)
	x1,y1:= player.steps[0].x,player.steps[0].y
	x2,y2:= player.steps[1].x,player.steps[1].y
	var orders [5]int=[5]int{-1,1,0,-2,2}

	for i:= 0 ;i<5;i++  {
		tmpx:=orders[i]+x1
	    if tmpx< 0 || tmpx >= 15 {
	        continue
	    }
	    for j := 0;j<5; j++ {
			tmpy:=orders[j]+y1
	        if tmpy < 0 || tmpy >= 15 {
	            continue
	        }
			dy:=(y1-y2)*2
			dx:=(x1-x2)*2
			switch{
			case x1==x2:
				if dy+y1==tmpy && (x1+dy==tmpx || x1-dy==tmpx){
					continue
				}
			case y1==y2:
				if dx+x1==tmpx && (y1+dx==tmpy || y1-dx==tmpy){
					continue
				}
			default:
				if x1+dx==tmpx && y1+dy==tmpy{
					continue
				}
			}
	        if player.frame[tmpx][tmpy] == 0 {
				sts = append(sts, StepInfo{tmpx, tmpy, 1,false})
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
		return &StepInfo{x:7,y:7,bw:1,forbid:false}
	case 1:
		return player.getWhiteFormula(player.steps[0].x,player.steps[0].y)
	case 2:
		if player.IsFormula(){
			return player.getThirdStep()
		}
	}
	return nil
}
