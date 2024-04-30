package snake

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (

	// MARGIN between game board and screen boundary
	MARGIN = 0

	// Number of wall updates per second
	TPS = 5

	NORMAL_FONT_SIZE = 13
)

type Game struct {
	SnakeState SnakeState

	game_tick_cnt          uint64
	snake_tick_cnt         uint64
	last_pressed_direction int
}

func CreateGame(height, width int) *Game {
	snake := CreateSnake(height, width)
	return &Game{*snake, 0, 0, snake.Direction}
}

func (g *Game) RestartGame() {
	snake := CreateSnake(g.SnakeState.Height, g.SnakeState.Width)
	g.SnakeState = *snake
}

func (g *Game) Update() error {
	g.game_tick_cnt = (g.game_tick_cnt + 1) % 60

	// Remember the last key pressed
	var keys []ebiten.Key
	keys = inpututil.AppendJustPressedKeys(keys)
	for _, key := range keys {
		switch key {
		case ebiten.KeyUp:
			g.last_pressed_direction = UP
		case ebiten.KeyDown:
			g.last_pressed_direction = DOWN
		case ebiten.KeyLeft:
			g.last_pressed_direction = LEFT
		case ebiten.KeyRight:
			g.last_pressed_direction = RIGHT
		case ebiten.KeyR:
			// Press R to restart
			g.RestartGame()
		}
	}

	// move the snake 5 times every second
	if !g.SnakeState.GameOver && g.game_tick_cnt%(60/TPS) == 0 {
		g.snake_tick_cnt += 1
		// Update direction
		g.SnakeState.UpdateDirection(g.last_pressed_direction)
		// move snake one tick
		g.SnakeState.Tick()
	}
	return nil
}

func draw_game_info(screen *ebiten.Image, x int, y int, msg string) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, msg, &text.GoTextFace{
		Source: M_PLUS_FACE_SCOURCE,
		Size:   NORMAL_FONT_SIZE,
	}, op)
}

func (g *Game) get_cell_size(screen *ebiten.Image) (float64, float64) {
	width := float64(screen.Bounds().Dx()) / float64(g.SnakeState.Width)
	height := float64(screen.Bounds().Dy()) / float64(g.SnakeState.Height)
	return width, height
}

func (g *Game) Draw(screen *ebiten.Image) {
	cell_width, cell_height := g.get_cell_size(screen)

	// Draw boarder
	draw_boarder(screen, g.SnakeState.Width, g.SnakeState.Height, cell_width, cell_height)

	// Print game stat, score and ticks played
	score_str := fmt.Sprintf("Score: %5d", g.SnakeState.Score)
	tick_str := fmt.Sprintf("Time:  %5d", g.snake_tick_cnt/TPS)
	draw_game_info(screen, int(cell_width)+10, int(cell_height)+10, score_str)
	draw_game_info(screen, int(cell_width)+130, int(cell_height)+10, tick_str)

	if g.SnakeState.GameOver {
		draw_game_info(screen, screen.Bounds().Dx()/2-30, screen.Bounds().Dy()/2, "Game Over")
	}

	// Draw snake
	for _, snake_part := range g.SnakeState.SnakeBody {
		draw_snake_part(screen, snake_part, cell_width, cell_height)
	}

	// Draw apple
	if g.SnakeState.HasApple() {
		draw_apple(screen, g.SnakeState.Apple, cell_width, cell_height)
	}
}

// Draw game board boarder
func draw_boarder(screen *ebiten.Image, row, col int, cell_width, cell_height float64) {
	// First and last row
	vector.DrawFilledRect(
		screen,
		0, 0,
		float32(screen.Bounds().Dx()), float32(cell_height),
		BOARDER_COLOR, false)
	vector.DrawFilledRect(
		screen,
		0, float32(row-1)*float32(cell_height),
		float32(screen.Bounds().Dx()), float32(cell_height),
		BOARDER_COLOR, false)

	// first and last column
	vector.DrawFilledRect(
		screen,
		0, 0,
		float32(cell_width), float32(col)*float32(cell_height),
		BOARDER_COLOR, false)
	vector.DrawFilledRect(
		screen,
		float32(col-1)*float32(cell_width), 0,
		float32(cell_width), float32(row)*float32(cell_height),
		BOARDER_COLOR, false)

}

func draw_snake_part(screen *ebiten.Image, snake_part SnakePart, cell_width, cell_heigth float64) {
	body_part_img, ok := BODY_PART_TO_IMG_MAP[snake_part.PartType]
	if !ok {
		log.Fatalf("Unknow body type %v", snake_part)
	}
	screen.DrawImage(
		body_part_img,
		cell_img_option(
			snake_part.Cord.X,
			snake_part.Cord.Y,
			cell_width,
			cell_heigth))
}

func draw_apple(screen *ebiten.Image, apple_cord Point, cell_width, cell_height float64) {
	screen.DrawImage(
		APPLE_IMG,
		cell_img_option(
			apple_cord.X,
			apple_cord.Y,
			cell_width,
			cell_height))
}

// Returns the image option for drawing a cell
// The option will contain tranlate x, y and scale
func cell_img_option(col, row int, cell_width, cell_height float64) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	// Scale the image so that it fits in one cell
	op.GeoM.Scale(CellScale(cell_width), CellScale(cell_height))
	// Move the image to dst col and row
	op.GeoM.Translate(float64(col)*cell_width, float64(row)*cell_height)
	return op
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	s := min(outsideWidth, outsideHeight)
	return s - MARGIN/2, s - MARGIN/2
}
