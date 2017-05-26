package main

import (
	"ai"
	"fmt"
)

var l1,l2 int
/*
func Draw(player *ai.AIPlayer) {
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
			if j==x && i==y{
				bstr=" \033[7mx\033[0m"
				wstr=" \033[7mo\033[0m"
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
}*/

func simulate(show bool) int {
	player1, _ := ai.InitPlayer(1, l1, true)
	player2, _ := ai.InitPlayer(2, l2, true)
	over := 0
	for {
		x, y := player1.GetStep()
		if show {
			player1.Draw()
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
		}
		player2.SetStep(x, y)
		x, y = player2.GetStep()
		if show {
			player2.Draw()
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
	p.Draw()
	for {
		var x, y int
		fmt.Scanln(&x, &y)
		p.SetStep(x, y)
		p.Draw()
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
		}
		fmt.Printf("Total %d times, black win %d, white win %d, Drawn %d\n", color, bw, ww, dw)
		return
	}
	fmt.Println("Start:")
	player, err := ai.InitPlayer(color,2 , true)
	//player, err := ai.InitPlayer(color, 2, true)
	if err != nil {
		fmt.Println("Init server error:", err)
		return
	}
	over := 0
	if color == ai.BLACK {
		player.GetStep()
	}
	player.Draw()
	for ; over == 0; over = player.IsOver() {
		var x, y int
		fmt.Scanln(&x, &y)
		player.SetStep(x, y)
		player.Draw()
		if over = player.IsOver(); over != 0 {
			break
		}
		x, y = player.GetStep()
		player.Draw()
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
