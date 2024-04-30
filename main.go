package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/redwookcreek/snake/snake"
)

func main() {
	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Snake")
	if err := ebiten.RunGame(snake.CreateGame(20, 20)); err != nil {
		log.Fatal(err)
	}
}
