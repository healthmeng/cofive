package main

import (
"fmt"
"log"
"time"
"math/rand"
"ai"
"bufio"
"net"
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
	rand.Seed(time.Now().UnixNano())
	return rand.Int63()
}

var tkdata map[int] *ai.AIPlayer

func ReplyOver(over int){
}

// post setsteps, 
// if over return result/reset
// if not over (go getstep, store id) return situation and id
// post id get step and situation(and result)

func CreateBySteps(info *StepsInfo)(*ai.AIPlayer,int64){
	cur:=len(info.x)
	robot:=ai.WHITE
	id:=CreateRand()
	if cur%2 ==0{// next is black, so ai use black
		robot=ai.BLACK
	}
	p,_:=ai.InitPlayer(robot,info.level,info.forbid)
	for k,v:=range info.x{
		p.SetStep(v,info.y[k])
		if over:=p.IsOver();over!=0{
		//	ReplyOver(over)
			break
		}
	}
	return p,id
}

func ProcCurrent(conn net.Conn,p *ai.AIPlayer, id int64){
}

func ProcNext(conn net.Conn, id int64){
}

func procConn(conn net.Conn){
	defer conn.Close()
	rd:=bufio.NewReader(conn)
	line,_,err:=rd.ReadLine()
	if err!=nil{
		log.Println("Get command error:",err)
		return
	}
	switch string(line){
	case "PostStep":
		var steps StepsInfo
		line,_,err=rd.ReadLine()
		var num int
		fmt.Sscanf(string(line),"%d%d%d",&num,&steps.forbid,&steps.level)
		steps.x=make([]int,num)
		steps.y=make([]int,num)
		for i:=0;i<num;i++{
			line,_,err=rd.ReadLine()
			x,y:=0,0
			fmt.Sscanf(string(line),"%d%d",&x,&y)
		}
		p,id:=CreateBySteps(&steps)
		ProcCurrent(conn,p,id)
	case "GetFromID":
		line,_,err=rd.ReadLine()
		var id int64
		fmt.Sscanf(string(line),"%d",&id)
		ProcNext(conn,id)
	}
}

func main(){
	tkdata=make(map[int] *ai.AIPlayer)
	lis,err:=net.Listen("tcp",":547")
	if err!=nil{
		fmt.Println("Svr listen error:",err)
		return
	}
	defer lis.Close()
	for{
		if conn,err:=lis.Accept();err!=nil{
			fmt.Println("Accept error:",err)
		}else{
			go procConn(conn)
		}
	}
}

