package ai

/*
requirement:
1. common AI
2. concurrent computing
3. remote service
*/

import (
_"log"
"errors"
"time"
"math/rand"
)

const(
	MAX_STEP=15*15
	BLACK=1
	WHITE=2
	SCORE_INIT=-20000
)

type StepInfo struct{
	x,y int
	bw int
}

type AIPlayer struct{
	frame [15][15] int
	level int
	steps []StepInfo
	curstep int
	robot, human int
	rnd	*rand.Rand
}

func InitPlayer(color int, level int) (* AIPlayer,error){
	player:=new (AIPlayer)
	player.level=level
	player.steps=make([]StepInfo,MAX_STEP,MAX_STEP)
	player.rnd=rand.New(rand.NewSource(time.Now().UnixNano()))
	if player.robot=color;color==BLACK{
		player.human=WHITE
	}else if color==WHITE{
		player.human=BLACK
	}else{
		return nil,errors.New("Bad color")
	}

	return player, nil
}

func (player* AIPlayer)SetStep(x int,y int){
	player.frame[x][y]=player.human
	st:=StepInfo{x,y,player.human}
	player.steps[player.curstep]=st
	player.curstep++
}

func (player* AIPlayer)GetStep()(int,int){
	var st StepInfo
	if player.level==0{
		st=player.DirectAlgo()
	}else{
		st=player.MinMaxAlgo()
	}
	player.steps[player.curstep]=st
	player.curstep++
	return st.x,st.y
}

func (player* AIPlayer)ListSteps() ([]StepInfo,int){
	return player.steps,player.curstep
}

func (player* AIPlayer)Evaluate()(int,int){
	return 0,0
}

func (player* AIPlayer)chessaround(x,y int) bool{
	for i:=x-2;i<=x+2,i++{
		if x<0 || x>15{
			continue
		}
		for j:=y-2;j<y+2;j++{
			if y<0 || y>15{
				continue
			}
			if player.frame[i][j]{
				return true
			}
		}
	}
	return false
}

func (player* AIPlayer)getallstep()[]StepInfo{
	var sts:=make ([]StepInfo,0,MAX_STEP)
	for i:=0;i<15;i++{
		for j:=0;j<15;j++{
			if player.frame[i][j]==0 && chessaound(i,j) {
				sts=append(std,StepInfo{i,j,0})
			}
		}
	}
	return nil
}

func (player* AIPlayer)DirectAlgo()StepInfo{
	var st StepInfo
	allst:=player.getallstep()
	nstep:=len(allst)
//	if(nstep<1)
//		return nil
	results:=make ([]StepInfo,0,nstep)
	maxscore:=SCORE_INIT
	for i:=0;i<nstep;i++{
		player.ApplyStep(allst[i])
		bscore,wscore:=player.Evaluate()
		scores:=[2]int{bscore,wscore}
		val:=scores[player.computer]-scores[player.human]
		if val>maxscore{
			results=make([]stepInfo,1,nstep)
			results[0]=allst[i]
			maxscore=val
		}else if val==maxscore{
			results=append(results,allst[i])
		}
	}
	nchoose:=len(results)
	if nchoose>0{
		return results[rnd.Int()%nchoose]
	}
	return nil
}

func (player* AIPlayer)MinMaxAlgo()StepInfo{
	var st StepInfo
	return st
}

