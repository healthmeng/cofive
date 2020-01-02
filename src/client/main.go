package main

import (
	"ai"
	"errors"
	"fmt"
	"net"
	"bufio"
	"time"
)

var l1,l2 int

func simulate(show bool) (winner int,steps int, tm float64){
	player1, _ := ai.InitPlayer(1, l1, true)
	player2, _ := ai.InitPlayer(2, l2, true)
	over := 0
	tms:=time.Now()
	defer func (){
		d:=time.Since(tms)
		step:=player1.TotalSteps()
		if step<player2.TotalSteps(){
			step=player2.TotalSteps()
		}
		step-=3
		if step>0{
			fmt.Printf("Total %.1f seconds for %d steps, average %.2f seconds per step.\n",d.Seconds(),step,d.Seconds()/float64(step))
		}
		steps=step
		tm=d.Seconds()
	}()
	for {
		x, y := player1.GetStep(show)
		if show {
			player1.Draw(true)
			player1.DebugStep()
		}
		over = player1.IsOver()
		if over == 1 {
			if show{
				fmt.Println("Black win")
			}
			winner= 1
			return
		} else if over == -1 {
			if show{
			fmt.Println("Drawn...")
			}
			winner= -1
			return
		}else if over==2{
			if show{
			fmt.Println("White win")
			}
			winner=2
			return
		}
		player2.SetStep(x, y)
		x, y = player2.GetStep(show)
		if show {
			player2.Draw(true)
			player2.DebugStep()
		}
		over = player2.IsOver()
		if over != 0 {
			if show{
			fmt.Println("White win")
			}
			winner=2
			return
		} else if over == -1 {
			if show{
			fmt.Println("Drawn...")
			}
			winner= -1
			return
		}

		player1.SetStep(x, y)
	}
}

func manual() {
	//	p1:=ai.InitPlayer(BLACK)
	//	p2:=ai.InitPlayer(WHITE)
	p, _ := ai.InitPlayer(ai.WHITE, 0, true)
	p.Draw(true)
	for {
		var x, y int
		fmt.Scanf("%d%d",&x, &y)
		p.SetStep(x, y)
		p.Draw(true)
		if over := p.IsOver(); over != 0 {
			if over == 1 {
				fmt.Println("Black win")
			} else if over == 2 {
				fmt.Println("White win")
			} else {
				fmt.Println("Drawn")
			}
			break
		}
	}
}

func main() {
	fmt.Println("Robot use Black(1) or White(2)?")
	color := 0
	fmt.Scanln(&color)
	if color == 0 {
		manual()
		return
	}
	if color > 2 {
		bw, ww, dw := 0, 0, 0
		show := false
		if color%2 == 0 {
			show = true
		}
	fmt.Println("Player1,2 level:")
	fmt.Scanf("%d%d",&l1,&l2)
	totalstep:=0
	var totaltime float64=0.0
		for i := 0; i < color; i++ {
			winner,s,t:=simulate(show)
			totalstep+=s
			totaltime+=t
			switch winner {
			case 1:
				bw++
			case 2:
				ww++
			case -1:
				dw++
			}
			fmt.Println("Current win times: black/white/drawn",bw,ww,dw)
		}
		fmt.Printf("player1:%d, player2:%d. Total %d times, black win %d, white win %d, Drawn %d, total %d steps, %.1f secs, %.2f secs per step.\n",l1,l2, color, bw, ww, dw,totalstep,totaltime,totaltime/float64(totalstep))
		return
	}
	fmt.Println("AI level:")
	var al int
	fmt.Scanln(&al)
	fmt.Println("Start:")

	netcolor:=-1
	var player *ai.AIPlayer
	var err error
	if color<0{
		if color==-1{
			netcolor=ai.BLACK
		} else if color==-2{
			netcolor=ai.WHITE
		}
		player, err = ai.InitPlayer(netcolor,al, true)
	}else{
		player, err = ai.InitPlayer(color,al, true)
	}
	//player, err := ai.InitPlayer(color, 2, true)
	if err != nil {
		fmt.Println("Init server error:", err)
		return
	}
	over := 0
	if color == ai.BLACK || netcolor==ai.BLACK {
		player.GetStep(true)
	}
	player.Draw(true)
	var id int64 =-1
	for ; over == 0; /*over = player.IsOver() */{
		var x, y int
		fmt.Scanln(&x, &y)
		if x>=0 && y>=0 && x<15 && y<15{
			player.SetStep(x, y)
		}else{
			if x== -1 && y == -1{
				player.Retreat()
			}
			player.Draw(false)
			continue
		}
		player.Draw(true)
		var err error
		if color<0{
			id,over,err=SendNet(player)
			if err!=nil{
				fmt.Println("network error:",err)
				return
			}
		}else{
			over = player.IsOver()
		}
		if over!=0{
			break
		}
		if color<0{
			over,err=GetFromNet(player,id)
			if err!=nil{
				fmt.Println("network error:",err)
				return
			}
		}else{
			player.GetStep(true)
			over=player.IsOver()
		}
		player.Draw(true)
		player.DebugStep()
	}
	if over == 1 {
		fmt.Println("Black win")
	} else if over == 2 {
		fmt.Println("White win")
	} else {
		fmt.Println("Drawn")
	}
}

var rip string ="127.0.0.1"
var rport string=":547"

func SendNet(p *ai.AIPlayer)(int64,int,error){
	conn,err:=net.Dial("tcp",rip+rport)
	if err!=nil{
		return 0,0,err
	}
	defer conn.Close()
	steps,cnt:=p.ListSteps()
	sndstr:=fmt.Sprintf("PostStep\n%d %d %d\n",cnt,p.GetForbidInt(),p.GetAILevel())
	for _,pos:=range steps{
		sndstr+=fmt.Sprintf("%d %d\n",pos.X,pos.Y)
	}
	fmt.Println("sending:",sndstr)
	conn.Write([]byte(sndstr))
	rb:=bufio.NewReader(conn)
	line,_,_:=rb.ReadLine()
	var over,bval,wval int
	var id int64
	fmt.Sscanf(string(line),"%d%d%d%d",&over,&id,&bval,&wval)
	return id,over,nil
}

func GetFromNet(p *ai.AIPlayer,id int64)(int,error){
	conn,err:=net.Dial("tcp",rip+rport)
	if err!=nil{
		return 0,err
	}
	conn.Write([]byte(fmt.Sprintf("GetFromID\n%d\n",id)))
	rb:=bufio.NewReader(conn)
	line,_,_:=rb.ReadLine()
	if string(line)=="OK"{
		line,_,_=rb.ReadLine()
		var x,y,over,bval,wval int
		fmt.Sscanf(string(line),"%d%d%d%d%d",&x,&y,&over,&bval,&wval)
		p.SetStep(x,y)
		return over,nil
	}else{
		return 0,errors.New("id error")
	}
}
