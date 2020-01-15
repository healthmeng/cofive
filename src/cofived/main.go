package main

import (
"fmt"
"log"
"io/ioutil"
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

type GameInfo struct{
	Level int
	Num int
	Forbid int
	Steps string
}

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
/*		if err:=r.ParseForm();err!=nil{
			fmt.Println("parse form error:",err)
		}else{
//			fmt.Printf("level:%s, forbid %s, num:%s", r.Form["level"],r.Form["forbid"],r.Form["num"])
			for k,v:=range r.PostForm{
				fmt.Println("Key:",k,"Value:",v)
			}
		}*/
		result,_:=ioutil.ReadAll(r.Body)
		r.Body.Close()
		fmt.Println(string(result))
		var gminfo GameInfo
		if err:=json.Unmarshal(result,&gminfo);err!=nil{
			fmt.Println("Unmarshal error:",err)
			return
		}
		var steps StepsInfo
		steps.level=gminfo.Level
		if gminfo.Forbid==1{
			steps.forbid=true
		}
		steps.x=make([]int,gminfo.Num)
		steps.y=make([]int,gminfo.Num)
		xyinfo:=strings.Split(gminfo.Steps,";")
		for i:=0;i<gminfo.Num;i++{
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
			ProcCurrent(p,id)
		}
	}
}

func ConfigGame(w http.ResponseWriter, r *http.Request){
    if r.Method=="GET"{
        r.ParseForm()
        t,_:=template.ParseFiles("select.htm")
        t.Execute(w,nil)
    }
}
func HttpGetID(w http.ResponseWriter, r *http.Request){
	if(r.Method=="GET"){
		r.ParseForm()
		curinfo:=CurStepInfo{
			BVal:0,
			WVal:0,
			X:7,
			Y:7,
			IsOver:0,
			Id:r.Form.Get("id")	}
		id,_:=strconv.ParseInt(curinfo.Id,10,64)
	    if id>0{ // start:
			maplock.RLock()
			cdata,exists:=tkdata[id]
			maplock.RUnlock()
			if !exists{
				return
			}
			curinfo.X,curinfo.Y,curinfo.IsOver,curinfo.BVal,curinfo.WVal=<-cdata.chRet,<-cdata.chRet,<-cdata.chRet,<-cdata.chRet,<-cdata.chRet
		}
		if obj,err:=json.Marshal(&curinfo);err!=nil{
            fmt.Println("Marshal error")
            return
        }else{
            fmt.Fprint(w,string(obj))
		}
	}
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

	http.HandleFunc("/",ConfigGame)
	http.HandleFunc("/StartGame",StartGame)
	http.HandleFunc("/PostStep",HttpPostStep)
	http.HandleFunc("/GetFromID",HttpGetID)
	err=http.ListenAndServe(":7777",nil)
	if err!=nil{
		fmt.Println("error:",err)
	}
}

