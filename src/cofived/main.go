package main

import (
"fmt"
"log"
"time"
"math/rand"
"ai"
"bufio"
"net"
"sync"
"net/http"
"html/template"
"strconv"
"strings"
"encoding/json"
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

type CurStepInfo struct{
	BVal int
	WVal int
	X	int
	Y	int
	IsOver int
	Id	string
}

func CreateRand() int64{
	rand.Seed(time.Now().UnixNano())
	return rand.Int63()
}

type CalcData struct{
	player *ai.AIPlayer
	chRet chan int
}
/*
type GameInfo struct{
	Level string 
	Color string
	Forbid string
}
*/
var tkdata map[int64] *CalcData
var maplock sync.RWMutex

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
func ProcCurrent(p* ai.AIPlayer, id int64){
/*		bval,wval:=p.GetCurValues()
		over:=p.IsOver()
		res:=fmt.Sprintf("%d %d %d %d\n",over,id,bval,wval)
		conn.Write([]byte(res))
		if over!=0{
			return
		}*/
		cdata:=CalcData{p,make (chan int)}
		maplock.Lock()
		tkdata[id]=&cdata
		maplock.Unlock()
		go func (player *ai.AIPlayer){
			defer func (){
				maplock.Lock()
				delete(tkdata,id)
				maplock.Unlock()
			}()
			x,y:=player.GetStep(false)
			over:=player.IsOver()
			bval,wval:=player.GetCurValues()
			cdata.chRet<-x
			select {
			case <-time.After(time.Second*30):
			case cdata.chRet<-y:
				cdata.chRet<-over
				cdata.chRet<-bval
				cdata.chRet<-wval
			}
		}(p)
}

func ProcNext(conn net.Conn, id int64){
	if id==-1{ // start:
		conn.Write([]byte("OK\n7 7 0 0 0\n"))
		return
	}
	maplock.RLock()
	cdata,exists:=tkdata[id]
	maplock.RUnlock()
	if !exists{
		conn.Write([]byte("ERROR\n"))
		return
	}
	x,y,over,bval,wval:=<-cdata.chRet,<-cdata.chRet,<-cdata.chRet,<-cdata.chRet,<-cdata.chRet
	conn.Write([]byte(fmt.Sprintf("OK\n%d %d %d %d %d\n",x,y,over,bval,wval)))
}

func procConn(conn net.Conn){
	defer conn.Close()
	rd:=bufio.NewReader(conn)
	line,_,err:=rd.ReadLine()
	var steps StepsInfo
	if err!=nil{
		log.Println("Get command error:",err)
		return
	}
	switch string(line){
	case "PostStep":
		line,_,err=rd.ReadLine()
		var num ,forbid int
		fmt.Sscanf(string(line),"%d%d%d",&num,&forbid,&steps.level)
		if forbid==0{
			steps.forbid=false
		}else{
			steps.forbid=true
		}
		steps.x=make([]int,num)
		steps.y=make([]int,num)
		for i:=0;i<num;i++{
			line,_,err=rd.ReadLine()
			fmt.Sscanf(string(line),"%d%d",&steps.x[i],&steps.y[i])
		}
		p,id:=CreateBySteps(&steps)
		bval,wval:=p.GetCurValues()
		over:=p.IsOver()
		res:=fmt.Sprintf("%d %d %d %d\n",over,id,bval,wval)
		conn.Write([]byte(res))
		if over!=0{
			return
		}
		ProcCurrent(p,id)
	case "GetFromID":
		line,_,err=rd.ReadLine()
		var id int64
		fmt.Sscanf(string(line),"%d",&id)
		ProcNext(conn,id)
	}
}

func StartGame(w http.ResponseWriter, r *http.Request){
   // w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
    //w.Header().Add("Access-Control-Allow-Headers", "Content-Type") 
	if r.Method=="GET"{
		r.ParseForm()
		t,_:=template.ParseFiles("front.htm")
		data:=make(map[string]string)
		data["Level"]=r.Form.Get("level")
		data["Color"]=r.Form.Get("color")
		data["Forbid"]=r.Form.Get("forbid")
		t.Execute(w,data)
	}
}

func HttpPostStep(w http.ResponseWriter, r *http.Request){
    //w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
    //w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	if r.Method=="POST"{
		if err:=r.ParseForm();err!=nil{
			fmt.Println("parse form error:",err)
		}else{
			fmt.Printf("level:%s, forbid %s, num:%s", r.Form["level"],r.Form["forbid"],r.Form["num"])
		}
		var steps StepsInfo
		steps.level,_=strconv.Atoi(r.FormValue("level"))
		if r.FormValue("forbid")=="1"{
			steps.forbid=true
		}
		num,_:=strconv.Atoi(r.FormValue("num"))
		steps.x=make([]int,num)
		steps.y=make([]int,num)
		points:=r.FormValue("steps") // a,b;c,d;e,f
		xyinfo:=strings.Split(points,";")
		for i:=0;i<num;i++{
			pt:=strings.Split(xyinfo[i],",")
			steps.x[i],_=strconv.Atoi(pt[0])
			steps.y[i],_=strconv.Atoi(pt[1])
		}
		p,id:=CreateBySteps(&steps)
		var curinfo CurStepInfo
		curinfo.BVal,curinfo.WVal=p.GetCurValues()
		curinfo.X,curinfo.Y=p.GetLastStep()
		curinfo.IsOver=p.IsOver()
		curinfo.Id=fmt.Sprintf("%d",id)
		if obj,err:=json.Marshal(&curinfo);err!=nil{
			fmt.Println("Marshal error")
			return
		}else{
			fmt.Fprint(w,string(obj))
		}
	}
}

func HttpGetID(w http.ResponseWriter, r *http.Request){
}

func main(){
	tkdata=make(map[int64]*CalcData )
	lis,err:=net.Listen("tcp",":547")
	if err!=nil{
		fmt.Println("Svr listen error:",err)
		return
	}
	defer lis.Close()
	go func (){
		for{
			if conn,err:=lis.Accept();err!=nil{
				fmt.Println("Accept error:",err)
			}else{
				go procConn(conn)
			}
		}
	}()

	http.HandleFunc("/",StartGame)
	http.HandleFunc("/PostStep",HttpPostStep)
	http.HandleFunc("/GetFromID",HttpGetID)
	err=http.ListenAndServe(":7777",nil)
	if err!=nil{
		fmt.Println("error:",err)
	}
}

