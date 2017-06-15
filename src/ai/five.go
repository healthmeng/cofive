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
"strings"
"runtime"
"sync"
"fmt"
"math/rand"
)

const(
	MAX_STEP=15*15
	BLACK=1
	WHITE=2
	SCORE_INIT= -2000000
	FORBIDDEN = -500000
	WIN= 1000000
)

const(
	NONE=iota
	CCCCC
	CCCCCC
	CCCC_C
	CCC_CC
	CC_CCC
	NCC_CCCN
	CCCC
	NCCCC
	CCC
	NCCC
	CC
	NCC
	C
	NC
	END
)

var BScoreTB [15]int =[...]int{
	0,	// NONE
	WIN,	// CCCCC
	WIN,	// CCCCCC
	10000,	//CCCC_C
	1150,	// CCC_CC
	1020,	// CC_CCC
	1000,	// NCC_CCCN
	10000,	// CCCC
	1000,	// NCCCC
	900,	// CCC
	300,	// NCCC
	300,	// CC 
	20,		// NCC
	0,		// C
	0,		// NC
}

var FScoreTB [15] int=[...]int{
	0,	// NONE
	WIN,	// CCCCC
	WIN,	// CCCCCC
	50000,	// CCCC_C
	50150,	// CCC_CC
	50020,	// CC_CCC
	50000,	// NCC_CCCN
	50000,	// CCCC
	50000,	// NCCCC
	2000,	// CCC
	300,	// NCCC
	300,	// CC 
	30,		// NCC
	0,		// C
	0,		// NC
}

var IsWin bool =false
var ncpus int=0
var maxvlock sync.RWMutex
var steplock sync.Mutex
//var rnd *rand.Rand

func init(){
	if strings.ToLower(runtime.GOOS)=="windows"{
		IsWin=true
	}
	ncpus=runtime.NumCPU()
}

type StepInfo struct{
	x,y int
	bw int
	forbid bool
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

type AIPlayer struct{
	frame [15][15] int
	level int
	forbid bool
	steps []StepInfo
//	bvalues, wvalues []int
	bshapes,wshapes []map[int] int
//	bshapes, wshapes map[int]int
	curstep int
	robot, human int
}

func (dst* AIPlayer)Clone(src *AIPlayer){
	for i:=0;i<15;i++{
		for j:=0;j<15;j++{
			dst.frame[i][j]=src.frame[i][j]
		}
	}
// already cloned in InitPlayer
//	dst.level=src.level
//	dst.forbid=src.forbid
//	dst.robot=src.robot
//	dst.human=src.human
	dst.curstep=src.curstep

// slices and map should be created already
//	dst.steps=make([]StepInfo,MAX_STEP,MAX_STEP)
//	dst.bshapes=make([]map[int]int,MAX_STEP,MAX_STEP)
//	dst.wshapes=make([]map[int]int,MAX_STEP,MAX_STEP)
	for i:=0;i<src.curstep;i++{
		dst.steps[i]=src.steps[i]
		dst.bshapes[i]=make(map[int]int)
		dst.wshapes[i]=make(map[int]int)
		for j:=0;j<END;j++{
			dst.bshapes[i][j]=src.bshapes[i][j]
			dst.wshapes[i][j]=src.wshapes[i][j]
		}
	}
}

func (player* AIPlayer)GetLastStep()(x,y int){
	x,y= -1,-1
	if cur:=player.curstep;cur>0{
		x,y=player.steps[cur-1].x,player.steps[cur-1].y
	}
	return
}

func (player* AIPlayer)DebugStep(){
	n:=player.curstep-1
	if n<0{
		return
	}
	b,w:=player.GetCurValues()
	fmt.Printf("step %d,%d  value %d-%d\n",player.steps[n].x,player.steps[n].y,b,w)
}

func InitPlayer(color int, level int, forbid bool) (* AIPlayer,error){
	player:=new (AIPlayer)
	player.level=level
	player.forbid=forbid
	player.steps=make([]StepInfo,MAX_STEP,MAX_STEP)
	player.bshapes=make([]map[int]int,MAX_STEP,MAX_STEP)
	player.wshapes=make([]map[int]int,MAX_STEP,MAX_STEP)
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
	if player.curstep<8{
		return 0
	}

	if player.curstep>=MAX_STEP{
		log.Println("Drawned")
		return -1
	}

//	x,y,bw:=player.steps[player.curstep-1].x,player.steps[player.curstep-1].y,player.steps[player.curstep-1].bw

	if player.steps[player.curstep-1].bw==BLACK{
		if player.bshapes[player.curstep-1][CCCCC]>0{	// the rule: when get five, ignore forbid
			return BLACK
		}else{
			if !player.forbid{
				if player.bshapes[player.curstep-1][CCCCCC]>0{
					return BLACK
				}
			}else  if  player.steps[player.curstep-1].forbid==true{
				return WHITE
			}
		}
	}else{
		if player.wshapes[player.curstep-1][CCCCC]>0 || player.wshapes[player.curstep-1][CCCCCC]>0{
			return WHITE
		}
	}
/*
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
*/

	return 0
}

func (player* AIPlayer)SetStep(x int,y int){
	player.ApplyStep(StepInfo{x,y,player.curstep%2+1,false})
}

func (player* AIPlayer)ApplyStep(st StepInfo){
	var bshapes,wshapes map[int] int
	if player.curstep>0{
		bshapes,wshapes,_=player.CalShape(st.x,st.y,false)
	}
	curb:=make(map[int]int)
	curw:=make(map[int]int)
	player.frame[st.x][st.y]=st.bw
	player.steps[player.curstep]=st
	player.curstep++
	nbshapes,nwshapes,forbid:=player.CalShape(st.x,st.y,true)

	player.steps[player.curstep-1].forbid=forbid
	if player.curstep>1{
	// remove old
		for i:=1;i<END;i++{
			if player.bshapes[player.curstep-2][i]< bshapes[i]{
				log.Printf("Error! bshapes %d count :%d < %d\n",i,player.bshapes[player.curstep-2][i],bshapes[i])
			}
			curb[i]=player.bshapes[player.curstep-2][i]
			curw[i]=player.wshapes[player.curstep-2][i]
			if bshapes[i]!=0{
				curb[i]=player.bshapes[player.curstep-2][i]-bshapes[i]
			}
			if wshapes[i]!=0{
				curw[i]=player.wshapes[player.curstep-2][i]-wshapes[i]
			}
		}
	}

// add new
	for i:=1;i<END;i++{
		if nbshapes[i]!=0{
			curb[i]+=nbshapes[i]
		}
		if nwshapes[i]!=0{
			curw[i]+=nwshapes[i]
		}
	}
	player.bshapes[player.curstep-1]=curb
	player.wshapes[player.curstep-1]=curw
}

func (player* AIPlayer)UnsetStep(x,y int){
	player.frame[x][y]=0
	if player.curstep>0{
		player.curstep--
	}
}

func (player* AIPlayer)GetStep(debug bool)(int,int){
	var st *StepInfo=nil
	if player.curstep<3{
		st=player.TryFormula()
	}
	if st==nil{
		if player.level==0{
			st=player.DirectAlgo()
		}else{
			st=player.MinMaxAlgo(debug)
		}
		if st==nil{
		//	log.Println("Drawn!")
			return -1,-1
		}
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

func (part *Conti)ParseType()int{
	if part.conttype!=-1 {// already parsed before
		return part.conttype
	}
	cont:=part.length
	ret:=NONE
	midsp:=1
	if part.spmid==0{
		midsp=0
	}

	switch {
		case cont==1:
			if part.leftsp+part.rightsp>=4{
				if part.leftsp>0 && part.rightsp>0 && part.leftsp+part.rightsp>4 {
					ret=C
				}else{
					ret=NC
				}
			}
		case cont==2:
			if part.leftsp+part.rightsp+midsp>=3{
				if part.leftsp>0 && part.rightsp>0 && part.leftsp+part.rightsp+midsp>3{
					ret=CC
				}else{
					ret=NCC
				}
			}
		case cont==3:
			if part.leftsp+part.rightsp+midsp>=2{
				if part.leftsp>0 && part.rightsp>0 && part.leftsp+part.rightsp+midsp>2{
					ret=CCC
				}else{
					ret=NCCC
				}
			}
		case cont==4:
			if part.leftsp+part.rightsp+midsp>=1{
				if part.leftsp>0 && part.rightsp>0 && midsp==0{
					ret=CCCC
				}else{
					ret=NCCCC
				}
			}
		case cont>=5:
			maxcont:=part.spmid
			if maxcont<cont-part.spmid{
				maxcont=cont-part.spmid
			}
			if maxcont==5{
					ret=CCCCC
			}else if maxcont>5{
					ret=CCCCCC
			} else{
				max:=0
				if maxcont!=part.spmid{	// right part is longer
					if maxcont+part.rightsp>4{
						max=maxcont
					}else if part.spmid+part.leftsp>4{
						max=part.spmid
					}
				}else{ // left part
					if part.spmid+part.leftsp>4{
						max=part.spmid
					}else if maxcont+part.rightsp>4{
						max=maxcont
					}
				}
				switch max{
					case 4:
						ret=CCCC_C
					case 3:
						ret=CCC_CC
					case 2:
						ret=CC_CCC
					default:
						ret=NCC_CCCN
				}
			}
	}
	part.conttype=ret
	return ret
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

func (player* AIPlayer)GetCurValues()(int,int){
	if player.curstep<1{
		return 0,0
	}
/*	over:=player.IsOver()
	if over==BLACK{// check black in CheckForbid already
		return WIN,0
	}else if over==WHITE{
		return 0,WIN
	}*/
	nextmove:=player.curstep%2+1
	bval,wval:=0,0
	bd3,wd3,b4,w4:=0,0,0,0
	var btable, wtable []int
	if nextmove==BLACK{
		btable=FScoreTB[:]
		wtable=BScoreTB[:]
	}else{
		btable=BScoreTB[:]
		wtable=FScoreTB[:]
	}
	if player.forbid{
		btable[CCCCCC]= -WIN
		btable[CCCC_C]= btable[NCCCC]
		btable[CCC_CC]= btable[NCCC]
		btable[CC_CCC]= btable[NCC]
		btable[NCC_CCCN]=0
	}

	bnc3:=0
bout:
	for k,v:= range player.bshapes[player.curstep-1]{
		switch k{
		case CCCCC :
			if v>0{
				bval=WIN
				break bout
			}
		case CCCCCC :
			if v>0{
				if player.forbid{
					bval= -WIN
				}else{
					bval=WIN
				}
				break bout
			}
		case CCC:
			bd3=v
		case NCCCC :
			fallthrough
		case CCCC_C:
			fallthrough
		case CCCC:
			b4+=v
		case NCCC:
			bnc3=v
		}
		bval+=btable[k]*v
		if bnc3>1{
			bval+=(bnc3-1)*300
		}
	}
	wnc3:=0
wout:
	for k,v:= range player.wshapes[player.curstep-1]{
		switch k{
		case CCCCC :
			fallthrough
		case CCCCCC :
			if v>0{
				wval=WIN
				break wout
			}
		case CCC:
			wd3=v
		case NCCCC :
			fallthrough
		case CCCC_C :
			fallthrough
		case CCCC :
			w4+=v
		case NCCC:
			wnc3=v
		}
		wval+=wtable[k]*v
		if wnc3>1{
			wval+=(wnc3-1)*300
		}
	}
	if nextmove==WHITE{// forbid is excluded in IsOver()
		if bd3>=1 && b4>=1 && w4<1{	// 4-3
			bval+=10000
		}
		if b4>1 && w4<1{	// 4 occurred in 2 places
			bval+=10000
		}else if bd3 >1 && (w4<1 && wd3<1){	// double 3
			bval+=5000
		}
	}else if nextmove==BLACK{
		if wd3>=1 && w4>=1 && b4<1{
			wval+=10000
		}
		if w4>1 && b4<1{
			wval+=10000
		}
		if wd3>1 && (b4<1 && bd3<1){
			wval+=5000
		}
	}
	return bval,wval
}

func (player* AIPlayer) CountShape(parts []Conti,checkfb bool)(map [int]int,map[int]int,int/* forbid type */){
	bs:=make(map[int]int)
	ws:=make(map[int]int)

	nCCC:=0
	nCCCC:=0
	ftype:=0
	for  _,part:=range parts{
		tp:=part.ParseType()
		if part.bw==BLACK{
			bs[part.ParseType()]++
			if checkfb && player.forbid && ftype==0 && part.isnew{
				switch tp{
					case CCC:
						nCCC++
						if nCCC>=2{
							ftype=CCC
						}
					case CCCC:
						fallthrough
					case CCCC_C:
						fallthrough
					case CCC_CC:
						fallthrough
					case CC_CCC:
						fallthrough
					case NCC_CCCN:
						fallthrough
					case NCCCC:
						nCCCC++
						if nCCCC>=2{
							ftype=CCCC
						}
					case CCCCCC:
						ftype=CCCCCC
					}
			}
		}else{
			ws[part.ParseType()]++
		}
	}
	return bs,ws,ftype
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
					if line[i+1]==end.bw{// **-**-
						end.AddTail(0)
					}else{
						end=nil
					}
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
					}else if rend== -1{
						if rfr!= -1{
							log.Println("Error: rend==-1 && rfr!= -1")
							return nil
						}
						parts=append(parts,*front)
						parts=append(parts,*end)
						front=&Conti{end.rightsp,0,0,1,line[i],-1,false}
						end=nil
					}
				}else{	// end==nil && front!=nil
					if rfr==1{
						end=&Conti{1,0,0,1,line[i],-1,false}
					}else if rfr== -1{
						parts=append(parts,*front)
						front=&Conti{front.rightsp,0,0,1,line[i],-1,false}
					}
				}
			}else{ // front==nil
				front=&Conti{contsp,0,0,1,line[i],-1,false}
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
}
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
			case CCCC_C:
				fallthrough
			case CCC_CC:
				fallthrough
			case CC_CCC:
				fallthrough
			case NCC_CCCN:
				fallthrough
			case NCCCC:
				nCCCC++
				if nCCCC>=2{
					return CCCC
				}
			case CCCCCC:
				return CCCCCC
			}
		}
	}
	return 0
}

func (player* AIPlayer)CheckForbid(x,y int) int{
	if player.forbid{
		lines,places:=player.CrossLines(x,y)
		parts:=make([]Conti,0,MAX_STEP)
		for i:=0;i<4;i++{
			parts=append(parts,player.CountLineParts(lines[i],places[i])...)
		}
		return player.hasforbid(parts)
	}

	return 0
}
*/
func (player* AIPlayer)CrossLines(x,y int)([][]int,[]int){
	hor:=make([]int,0,15)
	ver:=make([]int,0,15)
	topleft:=make([]int,0,15)
	topright:=make([]int,0,15)

	places:=make([]int,4,4) // place is the index of the new piece in 1-D line

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

func (player* AIPlayer)CalShape(x,y int, checkfb bool)(map[int]int,map[int]int,bool){
	lines,places:=player.CrossLines(x,y)

	parts:=make([]Conti,0,MAX_STEP)
	for i:=0;i<4;i++{
		parts=append(parts,player.CountLineParts(lines[i],places[i])...)
	}
	if len(parts)==0{
		return nil,nil,false
	}

	hasforbid:=false
/*	if player.forbid && player.hasforbid(parts)!=0{
		hasforbid=true
	}
*/
	bs,ws,hf:=player.CountShape(parts,checkfb)
	if hf!=0{
		hasforbid=true
	}
	return bs,ws,hasforbid
}

type PT struct {
    x, y int
}

func (player *AIPlayer) getallstep(side int) []StepInfo {
    sts := make([]StepInfo, 0, MAX_STEP)
    if player.curstep == 0 {
        if side != BLACK {
            log.Println("Error, first step should be black turn")
        }
        sts = append(sts, StepInfo{7, 7, side,false})
    } else {
        stmap := make(map[PT]bool)
		var orders [5]int=[5]int{-1,1,0,-2,2}
        for nst := player.curstep - 1; nst >= 0;nst-- { /*
                pts:=player.chessaround(player.steps[i].x,player.steps[i].y)
                for _,v:=pts*/
            x:= player.steps[nst].x
            y:= player.steps[nst].y
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
                        pt := PT{tmpx, tmpy}
                        if !stmap[pt] {
                            stmap[pt] = true
                            sts = append(sts, StepInfo{tmpx, tmpy, side,false})
                        }
                    }
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
	val:=0
	for i:=0;i<nstep;i++{
		player.ApplyStep(allst[i])
        over:=player.IsOver()
        if over== player.robot{
           player.UnapplyStep(allst[i])
			return &allst[i]
		}else if over==0 {
			bscore,wscore:=player.GetCurValues()//player.bvalues[player.curstep-1],player.wvalues[player.curstep-1]
			scores:=[2]int{bscore,wscore}
			val=scores[player.robot-1]-scores[player.human-1]
		}else if over==player.human{
			val= -WIN
		}
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
		rnd:=rand.New(rand.NewSource(time.Now().UnixNano()))
		return &results[rnd.Int()%nchoose]
	}
	return nil
}

func (player* AIPlayer)Draw(hl bool){
	frame := player.GetFrame()
	fmt.Print("  ")
	for i := 0; i < 15; i++ {
		fmt.Printf("%2d", i)
	}
	fmt.Println("")

	x,y:=player.GetLastStep()
	for i := 0; i < 15; i++ {
		fmt.Printf("%-2d", i)
		for j := 0; j < 15; j++ {
			bstr:=" x"
			wstr:=" o"
			if hl && j==x && i==y{
				if IsWin{
					bstr=" X"
					wstr=" O"
				}else{
					bstr=" \033[7mx\033[0m"
					wstr=" \033[7mo\033[0m"
				}
			}
			switch frame[j][i] {
			case 0:
				fmt.Printf(" .")
			case 1:
				fmt.Printf(bstr)
			case 2:
				fmt.Printf(wstr)
			default:
				fmt.Printf(" ?")
			}
		}
		fmt.Println("")
	}
}

func (player* AIPlayer)TraceAll(){
	fmt.Println("\n\nStart Trace...")
	player.Draw(false)
	for i:=1;i<END;i++{
		fmt.Printf("bshape [%d]=%d, wshape[%d]=%d\n",
		i,
		player.bshapes[player.curstep-1][i],
		i,
		player.wshapes[player.curstep-1][i])
	}
	fmt.Println("-----End Trace-----")
}

func (player* AIPlayer)GetMax(x,y int,level int,beta int) int{
	if level==0{
		b,w:=player.GetCurValues()

		if player.robot==BLACK{
			return b-w
		}else{
			return w-b
		}
	}

	allst:=player.getallstep(player.robot) // always player.robot
	nstep:=len(allst)
	alpha:= SCORE_INIT
	if nstep<1{// no place left
		return 0 // drawn 
	}else{
		for i:=0;i<nstep;i++{
			var value int
			player.ApplyStep(allst[i])
			over:=player.IsOver()
			if over== player.robot{
				b,w:=player.GetCurValues()
				player.UnapplyStep(allst[i])
				if player.robot==BLACK{
					return b-w
				}else{
					return w-b
				}
			}else if over==player.human{
				value= -WIN
			}else{
				value=player.GetMin(allst[i].x,allst[i].y,level-1,&alpha)
			}
			if value>alpha{
				alpha=value
			}
			player.UnapplyStep(allst[i])
			if beta<= -WIN || alpha>=beta{
				break
			}
		}
	}
	return alpha
}

func (player* AIPlayer)GetMin(x,y int,level int, alpha *int) int{
	if level==0{
		b,w:=player.GetCurValues()
		if player.robot==BLACK{
			return b-w
		}else{
			return w-b
		}
	}

	allst:=player.getallstep(player.human) // always player.human
	nstep:=len(allst)
	beta:= -SCORE_INIT
	if nstep<1{// no place left
		return 0 // drawn 
	}else{
		for i:=0;i<nstep;i++{
			var value int
			player.ApplyStep(allst[i])
			over:=player.IsOver()
			if over== player.human{
				b,w:=player.GetCurValues()
				player.UnapplyStep(allst[i])
				if player.robot==BLACK{
					return b-w
				}else{
					return w-b
				}
			}else if over==player.robot{
				//log.Println("Error,impossible: black win in white turn")
				value=WIN	// forbidden
			}else{
				value=player.GetMax(allst[i].x,allst[i].y,level-1,beta)
			}
			if value<beta{
				beta=value
			}

			player.UnapplyStep(allst[i])
			if level==player.level-1{
				maxvlock.RLock()
				if *alpha>=WIN || beta<*alpha {
					maxvlock.RUnlock()
					break
				}
				maxvlock.RUnlock()
			}else{
				if *alpha>=WIN || beta<=*alpha{ // need not lock
					break
				}
			}
		}
	}
	return beta
}

func SearchPara(player *AIPlayer,steps[]StepInfo,finished chan int, max *int, maxsts *[]StepInfo, maxplayers *[]*AIPlayer){
	nstep:=len(steps)
	result:=0
	for i:=0;i<nstep;i++{
		player.ApplyStep(steps[i])
		over:=player.IsOver()
		if over== player.robot{// -1(drawn) is impossible, because nstep==1 will not enter SearchPara
			player.UnapplyStep(steps[i])
			maxvlock.Lock()
			*max=WIN
			*maxsts=make([]StepInfo,1,MAX_STEP)
			(*maxsts)[0]=steps[i]
			*maxplayers=make([]*AIPlayer,1,MAX_STEP)
			(*maxplayers)[0]=player
			maxvlock.Unlock()
			finished<-1
			return
		}else{
			var value int
			if over==player.human{
				value= -WIN
			}else{
				value=player.GetMin(steps[i].x,steps[i].y,player.level-1,max)
			}
			maxvlock.RLock()
			if *max<WIN && value>=*max{
				maxvlock.RUnlock()
				maxvlock.Lock()
				if value>*max{
					*maxsts=make([]StepInfo,1,MAX_STEP)
					(*maxsts)[0]=steps[i]
					*maxplayers=make([]*AIPlayer,1,MAX_STEP)
					(*maxplayers)[0]=player
					*max=value
					result=2 // new max
				} else if value==*max{
					*maxsts=append(*maxsts,steps[i])
					*maxplayers=append(*maxplayers,player)
					if result==0{
						result=3	// same as old max
					}
				}
				maxvlock.Unlock()
			}else{
				maxvlock.RUnlock()
			}
			player.UnapplyStep(steps[i])
		}
	}
	finished<-result
}

func (player* AIPlayer)MinMaxAlgo(debug bool ) *StepInfo{
	allst:=player.getallstep(player.robot) // always player.robot
	nstep:=len(allst)
	max:=SCORE_INIT
	maxsts:=make([]StepInfo,0,MAX_STEP)
	maxplayers:=make([]*AIPlayer,0,MAX_STEP)
	if nstep<1{
		return nil
	}else if nstep==1{
		return &allst[0]
	}else{
		finished:=make(chan int,ncpus)
		i:=0
		nslide:=nstep/ncpus
		if nslide==0{
			for i=0;i<nstep;i++{
				np,_:=InitPlayer(player.robot,player.level,player.forbid)
				np.Clone(player)
				go SearchPara(np,allst[i:i+1],finished,&max,&maxsts,&maxplayers)
			}
			for i=0;i<nstep;i++{
				<-finished
			}
		}else{
			for i=0;i<ncpus;i++{
				np,_:=InitPlayer(player.robot,player.level,player.forbid)
				np.Clone(player)
				end:=(i+1)*nslide
				if i==ncpus-1{
					end=nstep
				}
				go SearchPara(np,allst[i*nslide:end],finished,&max,&maxsts,&maxplayers)
			}
			for i=0;i<ncpus;i++{
				<-finished
			}
		}
		/*
			if i<nstep{
				if ithread<ncpus{
					// clone AIPlayer
					np,_:=InitPlayer(player.robot,player.level,player.forbid)
					np.Clone(player)
					ithread++
					i++
					go SearchPara(np,allst[i-1],finished,&max,&maxsts,&maxplayers)
				}else{
					ret:=<-finished
					ithread--
					if ret== 1{	// WIN OR DRAWN, the only step is steup
						for ithread>0{
							ret=<-finished
							ithread--
						}
						break
					}
					if ithread<0{
						log.Println("Error!, ithread<0")
					}
				}
			}else{
				for ithread>0{
					<-finished
					ithread--
				}
				break
			}*/
			}
/*		}else{
		}
	}*/
	nsts:=len(maxsts)
	if debug{
		fmt.Printf("%d calc result: %d step later: %d\n",player.robot,player.level,max)
	}
	if nsts>0{
		rnd:=rand.New(rand.NewSource(time.Now().UnixNano()))
		retindex:=rnd.Int()%nsts
		player.Clone(maxplayers[retindex])
		return &maxsts[retindex]
	}
	return nil
}

