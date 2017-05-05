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

var BScoreTB [11]int =[...]int{
	0,	// NONE
	50000,	// CCCCC
	900,	//CC_CCC
	5000,	// CCCC
	800,	// NCCCC,
	800,	// CCC
	150,	// NCCC
	200,	// CC 
	20,		// NCC
	20,		// C
	10,		// NC
}

var FScoreTB [11] int=[...]int{
	0,	// NONE
	50000,	// CCCCC
	15000,	//CC_CCC
	20000,	// CCCC
	15000,	// NCCCC
	1500,	// CCC
	200,	// NCCC
	200,	// CC 
	20,		// NCC
	20,		// C
	10,		// NC
}

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

func (player* AIPlayer) IsOver() int{
// 1 black, 2 white, 0 not over yet
	if player.curstep<9{
		return 0
	}
	x,y,bw:=player.steps[player.curstep-1].x,player.steps[player.curstep-1].y,player.steps[player.curstep-1].bw

	n:=0
	for i:=x-1;i>=0 && player.frame[i][y]==bw && n<4 ; i--{
		n++;
	}
	for i:=x+1;i<15 && player.frame[i][y]==bw && n<4; i++{
		n++;
	}
	if n>=4{
		return bw
	}

	n=0
	for i:=y-1 ;i>=0 && player.frame[x][i]==bw && n<4; i--{
		n++;
	}
	for i:=y+1;i<15 && player.frame[x][i]==bw && n<4; i++{
		n++;
	}
	if n>=4{
		return bw
	}

	n=0
	for i,j:=x-1,y-1;i>=0 && j>=0 && player.frame[i][j]==bw && n<4 ; i,j=i-1,j-1{
		n++;
	}
	for i,j:=x+1,y+1;i<15 && j<15 && player.frame[i][j]==bw && n<4;i,j=i+1,j+1{
		n++;
	}
	if n>=4{
		return bw
	}

	n=0
	for i,j:=x+1,y-1;i<15 && j>=0 && player.frame[i][j]==bw && n<4; i,j=i+1,j-1{
		n++;
	}
	for i,j:=x-1,y+1;i>=0 && j<15 && player.frame[i][j]==bw && n<4 ;i,j=i-1,j+1{
		n++;
	}
	if n>=4{
		return bw
	}

	return 0
}

func (player* AIPlayer)SetStep(x int,y int){
/*	bval,wval:=0,0
	if player.curstep>0{
		bval,wval=player.Evaluate(x,y)
	}
	player.frame[x][y]=player.human
	st:=StepInfo{x,y,player.human}
	player.steps[player.curstep]=st
	nbval,nwval:=player.Evaluate(x,y)
	if player.curstep>0{
		player.bvalues[player.curstep]= player.bvalues[player.curstep-1]+nbval-bval
		player.wvalues[player.curstep]= player.bvalues[player.curstep-1]+nwval-wval
	}else{
		player.bvalues[player.curstep]=nbval
		player.wvalues[player.curstep]=nwval
	}
	player.curstep++*/
	player.ApplyStep(StepInfo{x,y,player.human})
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
	if player.curstep>1{
		player.bvalues[player.curstep-1]= player.bvalues[player.curstep-2]+nbval-bval
		player.wvalues[player.curstep-1]= player.bvalues[player.curstep-2]+nwval-wval
	}else{
        player.bvalues[player.curstep-1]=nbval
        player.wvalues[player.curstep-1]=nwval
	}
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
	player.frame[st.x][st.y]=player.robot
	player.steps[player.curstep]=*st
	player.curstep++
	return st.x,st.y
}

func (player* AIPlayer)GetFrame() [15][15] int{
	return player.frame
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
	cont:=part.length%6
	ret:=NONE
	midsp:=1
	if part.spmid==0{
		midsp=0
	}

	switch cont{
		case 1:
			if part.leftsp+part.rightsp>4{
				if part.leftsp>0 && part.rightsp>0 {
					ret=C
				}else{
					ret=NC
				}
			}
		case 2:
			if part.leftsp+part.rightsp+midsp>3{
				if part.leftsp>0 && part.rightsp>0{
					ret=CC
				}else{
					ret=NCC
				}
			}
		case 3:
			if part.leftsp+part.rightsp+midsp>2{
				if part.leftsp>0 && part.rightsp>0{
					ret=CCC
				}else{
					ret=NCCC
				}
			}
		case 4:
			if part.leftsp+part.rightsp+midsp>1{
				if part.leftsp>0 && part.rightsp>0 && midsp==0{
					ret=CCCC
				}else{
					ret=NCCCC
				}
			}
		case 5:
			if midsp==0{
				ret=CCCCC
			}else{
				ret=CC_CCC
			}
	}
	part.conttype=ret
	return ret
}

func (part *Conti)CountScore(nextmove int) int{
// nextmove: 1-> black; 2->white; 0->ignore
	score:=0
	contype:=part.ParseType()
	if nextmove==part.bw{
		score=FScoreTB[contype]
	}else{
		score=BScoreTB[contype]
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
		if i<0 || i>=15{
			continue
		}
		for j:=y-2;j<y+2;j++{
			if j<0 || j>=15{
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
	if player.curstep==0{
		if side!=BLACK{
			log.Println("Error, first step should be black turn")
		}
		sts=append(sts,StepInfo{7,7,side})
	}else{
		for i:=0;i<15;i++{
			for j:=0;j<15;j++{
				if player.frame[i][j]==0 && player.chessaround(i,j) {
					sts=append(sts,StepInfo{i,j,side})
				}
			}
		}
	}
	return sts
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
		bscore,wscore:=player.bvalues[player.curstep-1],player.wvalues[player.curstep-1]
		scores:=[2]int{bscore,wscore}
		val:=scores[player.robot-1]-scores[player.human-1]
		log.Printf("%d,%d  value: %d-%d\n",allst[i].x,allst[i].y,bscore,wscore)
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
	log.Println("List:")
	for i:=0;i<nchoose;i++{
		log.Printf("%d,%d-%d\n",results[i].x,results[i].y,maxscore)
	}
	if nchoose>0{
		return &results[player.rnd.Int()%nchoose]
	}
	return nil
}

func (player* AIPlayer)MinMaxAlgo() *StepInfo{
	var st StepInfo
	return &st
}

