package main

import (
	"ai"
	"fmt"
)

var l1,l2 int

func simulate(show bool) int {
	player1, _ := ai.InitPlayer(1, l1, true)
	player2, _ := ai.InitPlayer(2, l2, true)
	over := 0
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
			return 1
		} else if over == -1 {
			if show{
			fmt.Println("Drawn...")
			}
			return -1
		}else if over==2{
			if show{
			fmt.Println("White win")
			}
			return 2
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
			return 2
		} else if over == -1 {
			if show{
			fmt.Println("Drawn...")
			}
			return -1
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
		fmt.Scanln(&x, &y)
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
		for i := 0; i < color; i++ {
			switch simulate(show) {
			case 1:
				bw++
			case 2:
				ww++
			case -1:
				dw++
			}
			fmt.Println("Current win times: black/white/drawn",bw,ww,dw)
		}
		fmt.Printf("player1:%d, player2:%d. Total %d times, black win %d, white win %d, Drawn %d\n",l1,l2, color, bw, ww, dw)
		return
	}
	fmt.Println("AI level:")
	var al int
	fmt.Scanln(&al)
	fmt.Println("Start:")
	player, err := ai.InitPlayer(color,al, true)
	//player, err := ai.InitPlayer(color, 2, true)
	if err != nil {
		fmt.Println("Init server error:", err)
		return
	}
	over := 0
	if color == ai.BLACK {
		player.GetStep(true)
	}
	player.Draw(true)
	for ; over == 0; over = player.IsOver() {
		var x, y int
		fmt.Scanln(&x, &y)
		if x>0 && y>0 && x<15 && y<15{
			player.SetStep(x, y)
		}else{
			if x== -1 && y == -1{
				player.Retreat()
			}
			player.Draw(false)
			continue
		}
		player.Draw(true)
		if over = player.IsOver(); over != 0 {
			break
		}
		x, y = player.GetStep(true)
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
