package main

import (
	"ai"
	"fmt"
)

func draw(player *ai.AIPlayer) {
	frame := player.GetFrame()
	fmt.Print("  ")
	for i := 0; i < 15; i++ {
		fmt.Printf("%2d", i)
	}
	fmt.Println("")
	for i := 0; i < 15; i++ {
		fmt.Printf("%-2d", i)
		for j := 0; j < 15; j++ {
			switch frame[j][i] {
			case 0:
				fmt.Printf(" .")
			case 1:
				fmt.Printf(" o")
			case 2:
				fmt.Printf(" x")
			default:
				fmt.Printf(" ?")
			}
		}
		fmt.Println("")
	}
}

func simulate(show bool) int {
	player1, _ := ai.InitPlayer(1, 0, true)
	player2, _ := ai.InitPlayer(2, 0, true)
	over := 0
	for {
		x, y := player1.GetStep()
		if show {
			draw(player1)
		}
		over = player1.IsOver()
		if over == 1 {
			fmt.Println("Black win")
			return 1
		} else if over == -1 {
			fmt.Println("Drawn...")
			return -1
		}
		player2.SetStep(x, y)
		x, y = player2.GetStep()
		if show {
			draw(player2)
		}
		over = player2.IsOver()
		if over != 0 {
			fmt.Println("White win")
			return 2
		} else if over == -1 {
			fmt.Println("Drawn...")
			return -1
		}

		player1.SetStep(x, y)
	}
}

func manual() {
	//	p1:=ai.InitPlayer(BLACK)
	//	p2:=ai.InitPlayer(WHITE)
	p, _ := ai.InitPlayer(ai.WHITE, 0, true)
	draw(p)
	for {
		var x, y int
		fmt.Scanln(&x, &y)
		p.SetStep(x, y)
		draw(p)
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
		fmt.Printf("Total %d times, black win %d, white win %d, drawn %d\n", color, bw, ww, dw)
		return
	}
	fmt.Println("Start:")
	player, err := ai.InitPlayer(color, 2, true)
	if err != nil {
		fmt.Println("Init server error:", err)
		return
	}
	over := 0
	if color == ai.BLACK {
		player.GetStep()
	}
	draw(player)
	for ; over == 0; over = player.IsOver() {
		var x, y int
		fmt.Scanln(&x, &y)
		player.SetStep(x, y)
		draw(player)
		if over = player.IsOver(); over != 0 {
			break
		}
		x, y = player.GetStep()
		draw(player)
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
