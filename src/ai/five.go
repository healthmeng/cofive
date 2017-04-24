package ai

/*
requirement:
1. common AI
2. concurrent computing
3. remote service
*/

import (
_"log"
)

const(
	MAX_STEP=15*15
	BLACK=1
	WHITE=2
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
}

func InitPlayer(color int, level int) (* AIPlayer,error){
	player:=new (AIPlayer)
	player.level=level
	player.steps=make([]StepInfo,MAX_STEP,MAX_STEP)
	return player, nil
}

func (player* AIPlayer)SetStep(x int,y int){
}

func (player* AIPlayer)GetStep()(x int,y int){
	return
}

