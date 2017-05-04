package ai

/*
requirement:
1. common AI
2. concurrent computing
3. remote service
*/

import (
"log"
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

const(
	NONE=iota
	CCCCC
	CC_CCC
	CCCC
	NCCCC
	CCC
	NCCC
	CC
	NCC
	C
	NC
)

var FScoreTB[11] int={0,50000,20000,}
var LScoreTB[11] int={0,50000,}

type StepInfo struct{
	x,y int
	bw int
}

type AIPlayer struct{
	frame [15][15] int
	level int
	steps []StepInfo
	bvalues, wvalues []int
	curstep int
	robot, human int
	rnd	*rand.Rand
}

func InitPlayer(color int, level int) (* AIPlayer,error){
	player:=new (AIPlayer)
	player.level=level
	player.steps=make([]StepInfo,MAX_STEP,MAX_STEP)
	player.bvalues=make([]int,MAX_STEP,MAX_STEP)
	player.wvalues=make([]int,MAX_STEP,MAX_STEP)
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
	bval,wval:=0,0
	if player.curstep>0{
		bval,wval=player.Evaluate(x,y)
	}
	player.frame[x][y]=player.human
	st:=StepInfo{x,y,player.human}
	player.steps[player.curstep]=st
	player.curstep++
	nbval,nwval:=player.Evaluate(x,y)
    player.bvalues[player.curstep]= player.bvalues[player.curstep-1]+nbval-bval
    player.wvalues[player.curstep]= player.bvalues[player.curstep-1]+nwval-wval
}

func (player* AIPlayer)UnsetStep(x,y int){
	player.frame[x][y]=0
	if player.curstep>0{
		player.curstep--
	}
}

func (player* AIPlayer)GetStep()(int,int){
	var st *StepInfo
	if player.level==0{
		st=player.DirectAlgo()
	}else{
		st=player.MinMaxAlgo()
	}
	player.steps[player.curstep]=*st
	player.curstep++
	return st.x,st.y
}

func (player* AIPlayer)ListSteps() ([]StepInfo,int){
	return player.steps,player.curstep
}

type Conti struct{
	leftsp int
	rightsp int
	spmid int
	length	int
	bw int // color
	conttype int
}

func (part *Conti)ParseType()int{
	return part.conttype
}

func (part *Conti)CountScore(nextmove int) int{
// nextmove: 1-> black; 2->white; 0->ignore
	score:=0
	maxcont:=part.length-part.spmid
	if maxcont<part.spmid{
		maxcont=part.spmid
	}
	switch {
	case part.length>=5:
		if maxcont>=5{
			score=50000
		}else{
			if nextmove==part.bw{
				score=20000
			}else{
				score=5000
			}
		}
	case part.length==4:
	/*	if maxcount==4){
			if nextmove==part.bw{
				score=20000
			}else{
				score=5000
			}
		}*/
	}
	return score
}

func (part *Conti)AddTail(cont int) int {
// return 0: if the part continue to be counted
// return 1: if part continue and a space in mid, need create end later
// return -1: if finished part, and a new part started
	var ret int
	switch {
		case cont==0:
			part.rightsp++
			ret=0
		case cont==part.bw:
			if part.rightsp==1{
				if part.spmid==0{ // **-*
					part.spmid=part.length
					part.length++
					part.rightsp=0
					ret=1
				}else{ // **-**-*
					ret= -1
				}
			}else if part.rightsp>1{
				ret= -1
			}else{ // rightsp==0
				part.length++
				ret=0
			}
		default:
			ret=-1
	}
	return ret
}

func (player* AIPlayer) CountScore(parts []Conti)(int,int){
	bs,ws:=0,0
	nextmove:=player.curstep%2+1
	for  _,part:=range parts{
		if part.bw==BLACK{
			bs+=part.CountScore(nextmove)
		}else{
			ws+=part.CountScore(nextmove)
		}
	}
	return bs,ws
}

func (player* AIPlayer)EvalLine(line[]int)(int,int){
	//bval,wval:=0,0
	nLen:=len(line)
	if nLen<5{
		return 0,0
	}
	parts:=make([]Conti,0,nLen)
	var front *Conti=nil
	var end *Conti=nil
	contsp:=0
// ...
	for i:=0;i<nLen;i++{
		switch {
		case i==nLen-1: // last, force check
			if front!=nil{
				front.AddTail(line[i])
				parts=append(parts,*front)
				if end!=nil{
					end.AddTail(line[i])
					parts=append(parts,*end)
				}
			}
		case line[i]==0:
			contsp++
			if front!=nil{
				front.AddTail(0)
				if end!=nil{
					end.AddTail(0)
				}
			}
		default: // line[i]!=0
			if front!=nil{
				rfr:=front.AddTail(line[i])
				if end!=nil{
					if rfr==1{
						log.Println("Error: rfr==1 && end!=nil")
						return 0,0
					}
				// rfr!=1
					rend:=end.AddTail(line[i])
					if rend==1{
						if rfr!= -1{
							log.Println("Error: rend==1 && rfr!=-1")
							return 0,0
						}
						// **-**-*
						parts=append(parts,*front)
						front=end
						end=&Conti{1,0,0,1,line[i],0}
					}else if rend== -1{
						if rfr!= -1{
							log.Println("Error: rend==-1 && rfr!= -1")
							return 0,0
						}
						parts=append(parts,*front)
						parts=append(parts,*end)
						front=&Conti{end.rightsp,0,0,1,line[i],0}
						end=nil
					}
				}else{	// end==nil && front!=nil
					if rfr==1{
						end=&Conti{1,0,0,1,line[i],0}
					}else if rfr==-1{
						parts=append(parts,*front)
						front=&Conti{front.rightsp,0,0,1,line[1],0}
					}
				}
			}else{ // front==nil
				front=&Conti{contsp,0,0,1,line[i],0}
			}
			contsp=0
		}
	}
	return player.CountScore(parts)
}

func (player* AIPlayer)Evaluate(x,y int)(int,int){
	hor:=make([]int,0,15)
	ver:=make([]int,0,15)
	topleft:=make([]int,0,15)
	topright:=make([]int,0,15)

	for i,j:=x,y;i<15 && j<15; i,j=i+1,j+1{
		hor=append(hor,player.frame[i][y])
		ver=append(ver,player.frame[x][i])
	}

	for i,j:=x,y;i>=0 && j>=0 ;i,j=i-1,j-1{
		topleft=append(topleft,player.frame[i][j])
	}
	for i,j:=x+1,y+1;i<15 && j<15;i,j=i+1,j+1{
		topleft=append(topleft,player.frame[i][j])
	}

	for i,j:=x,y;i<15 && j>=0;i,j=i+1,j-1{
		topright=append(topleft,player.frame[i][j])
	}
	for i,j:=x-1,y+1;i>=0 && j<15;i,j=i-1,j+1{
		topright=append(topleft,player.frame[i][j])
	}

	bvalue,wvalue:=0,0
	btmp,wtmp:=player.EvalLine(hor)
	bvalue+=btmp
	wvalue+=wtmp

	btmp,wtmp=player.EvalLine(ver)
	bvalue+=btmp
	wvalue+=wtmp

	btmp,wtmp=player.EvalLine(topleft)
	bvalue+=btmp
	wvalue+=wtmp

	btmp,wtmp=player.EvalLine(topright)
	bvalue+=btmp
	wvalue+=wtmp

	return bvalue,wvalue
}

func (player* AIPlayer)chessaround(x,y int) bool{
	for i:=x-2;i<=x+2;i++{
		if x<0 || x>15{
			continue
		}
		for j:=y-2;j<y+2;j++{
			if y<0 || y>15{
				continue
			}
			if player.frame[i][j]!=0{
				return true
			}
		}
	}
	return false
}

func (player* AIPlayer)getallstep(side int)[]StepInfo{
	sts:=make ([]StepInfo,0,MAX_STEP)
	for i:=0;i<15;i++{
		for j:=0;j<15;j++{
			if player.frame[i][j]==0 && player.chessaround(i,j) {
				sts=append(sts,StepInfo{i,j,side})
			}
		}
	}
	return nil
}

func (player* AIPlayer)ApplyStep(st StepInfo){
	bval,wval:=0,0
	if player.curstep>0{
		bval,wval=player.Evaluate(st.x,st.y)
	}
	player.frame[st.x][st.y]=st.bw
	player.steps[player.curstep]=st
	player.curstep++
	nbval,nwval:=player.Evaluate(st.x,st.y)
	player.bvalues[player.curstep]= player.bvalues[player.curstep-1]+nbval-bval
	player.wvalues[player.curstep]= player.bvalues[player.curstep-1]+nwval-wval
}

func (player* AIPlayer)UnapplyStep(st StepInfo){
	player.frame[st.x][st.y]=0
	if player.curstep>0{
		player.curstep--
	}
}

func (player* AIPlayer)DirectAlgo()*StepInfo{
	allst:=player.getallstep(player.robot)
	nstep:=len(allst)
//	if(nstep<1)
//		return nil
	results:=make ([]StepInfo,0,nstep)
	maxscore:=SCORE_INIT
	for i:=0;i<nstep;i++{
		player.ApplyStep(allst[i])
		bscore,wscore:=player.Evaluate(allst[i].x,allst[i].y)
		scores:=[2]int{bscore,wscore}
		val:=scores[player.robot]-scores[player.human]
		if val>maxscore{
			results=make([]StepInfo,1,nstep)
			results[0]=allst[i]
			maxscore=val
		}else if val==maxscore{
			results=append(results,allst[i])
		}
		player.UnapplyStep(allst[i])
	}
	nchoose:=len(results)
	if nchoose>0{
		return &results[player.rnd.Int()%nchoose]
	}
	return nil
}

func (player* AIPlayer)MinMaxAlgo() *StepInfo{
	var st StepInfo
	return &st
}

