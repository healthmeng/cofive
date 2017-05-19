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
	SCORE_INIT= -2000000
	FORBIDDEN = -500000
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
	150,	// NCCC
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
	forbid bool
	steps []StepInfo
	bvalues, wvalues []int
	curstep int
	robot, human int
	rnd	*rand.Rand
}

func (player* AIPlayer)DebugStep(){
	n:=player.curstep-1
	if n<0{
		return
	}
	log.Printf("step %d,%d  value %d-%d\n",player.steps[n].x,player.steps[n].y,player.bvalues[n],player.wvalues[n])
}

func InitPlayer(color int, level int, forbid bool) (* AIPlayer,error){
	player:=new (AIPlayer)
	player.level=level
	player.forbid=forbid
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
	if player.curstep<5{
		return 0
	}
	x,y,bw:=player.steps[player.curstep-1].x,player.steps[player.curstep-1].y,player.steps[player.curstep-1].bw
	if bw==BLACK && player.CheckForbid(x,y)!=0{
		return WHITE
	}
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

	if player.curstep>=MAX_STEP{
		log.Println("Drawned")
		return -1
	}

	return 0
}

func (player* AIPlayer)SetStep(x int,y int){
	player.ApplyStep(StepInfo{x,y,player.curstep%2+1})
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
		player.wvalues[player.curstep-1]= player.wvalues[player.curstep-2]+nwval-wval
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
	if st==nil{
		log.Println("Drawn!")
		return -1,-1
	}
	player.ApplyStep(*st)
//	log.Printf("x,y: %d-%d, val: %d...%d\n",st.x,st.y,player.bvalues[player.curstep-1],player.wvalues[player.curstep-1])
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
	isnew bool // created by latest chess
}

func (part *Conti)ParseType()int{
	if part.conttype!=-1 {// already parsed before
		return part.conttype
	}
	cont:=part.length%6
	ret:=NONE
	midsp:=1
	if part.spmid==0{
		midsp=0
	}

	switch cont{
		case 1:
			if part.leftsp+part.rightsp>=4{
				if part.leftsp>0 && part.rightsp>0 && part.leftsp+part.rightsp>4 {
					ret=C
				}else{
					ret=NC
				}
			}
		case 2:
			if part.leftsp+part.rightsp+midsp>=3{
				if part.leftsp>0 && part.rightsp>0 && part.leftsp+part.rightsp+midsp>3{
					ret=CC
				}else{
					ret=NCC
				}
			}
		case 3:
			if part.leftsp+part.rightsp+midsp>=2{
				if part.leftsp>0 && part.rightsp>0 && part.leftsp+part.rightsp+midsp>2{
					ret=CCC
				}else{
					ret=NCCC
				}
			}
		case 4:
			if part.leftsp+part.rightsp+midsp>=1{
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
// return 1: if part continue and a space in mid, need create "end" later
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
//	nextmove:=player.curstep%2+1
	for  _,part:=range parts{
		if part.bw==BLACK{
			bs+=part.CountScore(player.human)
		}else{
			ws+=part.CountScore(player.human)
		}
	}
	return bs,ws
}

func (player* AIPlayer)CountLineParts(line[]int,newone int)[]Conti{
	nLen:=len(line)
	if nLen<5{
		return nil
	}
	parts:=make([]Conti,0,nLen)
	var front *Conti=nil
	var end *Conti=nil
	contsp:=0

	for i:=0;i<nLen;i++{
		switch {
		case i==nLen-1: // last, force check
			if front!=nil{
				front.AddTail(line[i])
				if newone==i{
					front.isnew=true
				}
				parts=append(parts,*front)
				if end!=nil{
					end.AddTail(line[i])
					if newone==i{
						end.isnew=true
					}
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
						return nil
					}
				// rfr!=1
					rend:=end.AddTail(line[i])
					if rend==1{
						if rfr!= -1{
							log.Println("Error: rend==1(end create another end) && rfr!= -1")
							return nil
						}
						// **-**-*
						parts=append(parts,*front)
						front=end
						end=&Conti{1,0,0,1,line[i],-1,false}
				/*		if i==newone{
							front.isnew=true
							end.isnew=true
						}*/
					}else if rend== -1{
						if rfr!= -1{
							log.Println("Error: rend==-1 && rfr!= -1")
							return nil
						}
						parts=append(parts,*front)
						parts=append(parts,*end)
						front=&Conti{end.rightsp,0,0,1,line[i],-1,false}
						end=nil
					/*	if i==newone{
							front.isnew=true
						}*/
					}
				}else{	// end==nil && front!=nil
					if rfr==1{
						end=&Conti{1,0,0,1,line[i],-1,false}
               /*      if i==newone{
                            front.isnew=true
                            end.isnew=true
                        }*/
					}else if rfr== -1{
						parts=append(parts,*front)
						front=&Conti{front.rightsp,0,0,1,line[i],-1,false}
				/*		if i==newone{
                            front.isnew=true
                        }*/
					}
				}
			}else{ // front==nil
				front=&Conti{contsp,0,0,1,line[i],-1,false}
		/*		if newone==i{
					front.isnew=true
				}*/
			}
			if newone == i{
				if front!=nil{
					front.isnew=true
				}
				if end!=nil{
					end.isnew=true
				}
			}
			contsp=0
		}
	}
	return parts
}

/*
func (player* AIPlayer)EvalLine(line[]int, place int)(int,int){
	//bval,wval:=0,0
	parts:=player.CountParts(line,place)
	if parts==nil{
		return 0,0
	}
	return player.CountScore(parts)
}*/

func (player* AIPlayer)hasforbid(parts []Conti) int{
	if player.forbid==false || len(parts)==0{
			return 0
	}
	nCCC:=0
	nCCCC:=0
	for _,p:=range parts{

		if p.isnew{
			tp:=p.ParseType()
			switch tp{
			case CCC:
				nCCC++
				if nCCC>=2{
					return CCC
				}
			case CCCC:
				fallthrough
			case NCCCC:
				nCCCC++
				if nCCCC>=2{
					return CCCC
				}
			case CCCCC:
				if p.length>5{
					return CCCCC
				}
			}
		}
	}
	return 0
}

func (player* AIPlayer)CheckForbid(x,y int) int{
	if player.frame[x][y] == BLACK {
		lines,places:=player.CrossLines(x,y)
		parts:=make([]Conti,0,MAX_STEP)
		for i:=0;i<4;i++{
			parts=append(parts,player.CountLineParts(lines[i],places[i])...)
		}
		return player.hasforbid(parts)
	}

	return 0
}

func (player* AIPlayer)CrossLines(x,y int)([][]int,[]int){
	hor:=make([]int,0,15)
	ver:=make([]int,0,15)
	topleft:=make([]int,0,15)
	topright:=make([]int,0,15)

	places:=make([]int,4,4)

	for i:=0;i<15; i++{
		hor=append(hor,player.frame[i][y])
		ver=append(ver,player.frame[x][i])
	}
	places[0]=x
	places[1]=y

	for i,j:=x,y;i>=0 && j>=0 ;i,j=i-1,j-1{
		topleft=append(topleft,player.frame[i][j])
	}
	tmplen:=len(topleft)
	places[2]=tmplen-1
	for i:=0;i<tmplen/2;i++{
		topleft[i],topleft[tmplen-1-i]=topleft[tmplen-1-i],topleft[i]
	}
	for i,j:=x+1,y+1;i<15 && j<15;i,j=i+1,j+1{
		topleft=append(topleft,player.frame[i][j])
	}

	for i,j:=x,y;i<15 && j>=0;i,j=i+1,j-1{
		topright=append(topright,player.frame[i][j])
	}
    tmplen=len(topright)
	places[3]=tmplen-1
    for i:=0;i<tmplen/2;i++{
        topright[i],topright[tmplen-1-i]=topright[tmplen-1-i],topright[i]
    }
	for i,j:=x-1,y+1;i>=0 && j<15;i,j=i-1,j+1{
		topright=append(topright,player.frame[i][j])
	}

	lines:=make([][]int,0,4)
	lines=append(lines,hor,ver,topleft,topright)
	return lines,places
}

func (player* AIPlayer)Evaluate(x,y int)(int,int){
	lines,places:=player.CrossLines(x,y)

	parts:=make([]Conti,0,MAX_STEP)
	for i:=0;i<4;i++{
		parts=append(parts,player.CountLineParts(lines[i],places[i])...)
	}
	if len(parts)==0{
		return 0,0
	}

	if player.CheckForbid(x,y)!=0{// check black in CheckForbid already
		return 0,FORBIDDEN
	}

	return player.CountScore(parts)
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

func (player* AIPlayer)GetMax(x,y int,level int) int{
	if level==0{
		b,w:=player.Evaluate(x,y)
		if player.robot==BLACK{
			return b-w
		}else{
			return w-b
		}
	}
	allst:=player.getallstep(player.robot) // always player.robot
	nstep:=len(allst)
	min:= -SCORE_INIT
	if nstep<1{
        b,w:=player.Evaluate(x,y)
        if player.robot==BLACK{
            return b-w
        }else{
            return w-b
        }
	}else{
		for i:=0;i<nstep;i++{
			player.ApplyStep(allst[i])
			value:=player.GetMin(allst[i].x,allst[i].y,level-1)
			if value<min{
				min=value
			}
		}
		return min
	}
}

func (player* AIPlayer)GetMin(x,y int,level int) int{
	if level==0{
		b,w:=player.Evaluate(x,y)
		if player.robot==BLACK{
			return b-w
		}else{
			return w-b
		}
	}
	allst:=player.getallstep(player.robot) // always player.robot
	nstep:=len(allst)
	max:= SCORE_INIT
//	minsts:=make([]StepInfo,0,MAX_STEP)
	if nstep<1{
        b,w:=player.Evaluate(x,y)
        if player.robot==BLACK{
            return b-w
        }else{
            return w-b
        }

	}else{
		for i:=0;i<nstep;i++{
			player.ApplyStep(allst[i])
			value:=player.GetMax(allst[i].x,allst[i].y,level-1)
			if value>max{
				max=value
			}
		}
		return max
	}

}

func (player* AIPlayer)MinMaxAlgo(/*nlevel int should be even*/ ) *StepInfo{
	allst:=player.getallstep(player.robot) // always player.robot
	nstep:=len(allst)
	max:=SCORE_INIT
	maxsts:=make([]StepInfo,0,MAX_STEP)
	if nstep<1{
		return nil
	}else{
		for i:=0;i<nstep;i++{
			player.ApplyStep(allst[i])
			value:=player.GetMin(allst[i].x,allst[i].y,player.level)
			if value>max{
				maxsts=make([]StepInfo,0,MAX_STEP)
				maxsts=append(maxsts,allst[i])
				max=value
			}else if value==max{
				maxsts=append(maxsts,allst[i])
			}
		}
	}
	nsts:=len(maxsts)
	if nsts>0{
		return &maxsts[player.rnd.Int()%nsts]
	}
	return nil
}

