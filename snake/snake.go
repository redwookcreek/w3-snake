package snake

import (
	"fmt"
	"log"
	"math/rand"
)

type Point struct {
	X int
	Y int
}

// Part of a snake body, with a coordinate for the location
// and indicate for head/tail/body
type SnakePart struct {
	Cord     Point
	PartType int
}

// moving direction of the snake
const (
	UP    = iota
	LEFT  = iota
	DOWN  = iota
	RIGHT = iota
)

const (
	BODY_PART_HEAD_UP    = iota
	BODY_PART_HEAD_LEFT  = iota
	BODY_PART_HEAD_DOWN  = iota
	BODY_PART_HEAD_RIGHT = iota
	BODY_PART_TAIL_UP    = iota
	BODY_PART_TAIL_DOWN  = iota
	BODY_PART_TAIL_LEFT  = iota
	BODY_PART_TAIL_RIGHT = iota
	BODY_PART_I          = iota // a straight up/down body part
	BODY_PART_H          = iota // a straight left/right body part
	BODY_PART_BODY_L     = iota // corner unit of a body, shape like a L
	BODY_PART_BODY_L1    = iota // L turned clockwise 90 degree
	BODY_PART_BODY_L2    = iota // L turned clockwise 180 degree
	BODY_PART_BODY_L3    = iota // L turned clockwise 270 degree
)

// Represents the state of a snake game
type SnakeState struct {
	// snake body is a list of points
	// the SnakeBody[0] is the tail, SnakeBody[-1] is the head
	SnakeBody []SnakePart

	// Moving direction of the snake
	Direction int

	// Location of the apple
	Apple Point

	// size of the game board
	Height int
	Width  int

	// True if game over
	GameOver bool

	// Game score, number of apples ate
	Score int
}

func CreateSnake(height, width int) *SnakeState {
	snake_body := make([]SnakePart, 1)
	// Init the snake at center of the board
	snake_body[0] = SnakePart{
		Point{width / 2, height / 2},
		BODY_PART_HEAD_DOWN}

	return &SnakeState{
		snake_body,
		// snake is moving down at the start of the game
		DOWN,
		// -1, -1 represents no apple on board
		Point{-1, -1},

		// board size
		height,
		width,

		// game over
		false,

		// game score
		0,
	}
}

func (ss *SnakeState) HasApple() bool {
	return ss.Apple.X > 0 && ss.Apple.Y > 0
}

// Update the snake's direction
// The snake cannot reverse, e.g. changing from UP to DOWN
func (ss *SnakeState) UpdateDirection(dir int) {
	if ss.Direction == UP || ss.Direction == DOWN {
		// snake is moving up or down
		if dir == LEFT || dir == RIGHT {
			ss.Direction = dir
		}
	} else {
		// snake is moving left or right
		if dir == UP || dir == DOWN {
			ss.Direction = dir
		}
	}
}

// Create apple if does not exist
func (ss *SnakeState) maybe_create_apple() {
	if ss.Apple.X > 0 && ss.Apple.Y > 0 {
		// Apple already exist
		return
	}
	// row 0, height -1 and col 0, width -1
	// are saved for boarders
	x := rand.Intn(ss.Width-2) + 1
	y := rand.Intn(ss.Height-2) + 1
	ss.Apple = Point{x, y}
}

// Advance snake one tick
func (ss *SnakeState) Tick() {
	if ss.GameOver {
		return
	}
	new_head := ss.advance_snake_head()
	touched := ss.snake_touched(new_head.Cord)
	if touched {
		ss.GameOver = true
		return
	}

	ss.maybe_consume_apple_and_grow_snake(new_head)
	ss.maybe_create_apple()
}

// Advance snake head by one cell, return the new snake head
func (ss *SnakeState) advance_snake_head() SnakePart {
	head_idx := len(ss.SnakeBody) - 1
	old_head := ss.SnakeBody[head_idx]
	new_head := old_head
	new_head.PartType = _HEAD_TYPE_FROM_DIR[ss.Direction]
	switch ss.Direction {
	case UP:
		new_head.Cord.Y -= 1
	case DOWN:
		new_head.Cord.Y += 1
	case LEFT:
		new_head.Cord.X -= 1
	case RIGHT:
		new_head.Cord.X += 1
	}
	return new_head
}

// Returns true if new head touches boarder or snake itself
func (ss *SnakeState) snake_touched(new_head Point) bool {
	if new_head.X <= 0 || new_head.X >= ss.Width-1 {
		return true
	}
	if new_head.Y <= 0 || new_head.Y >= ss.Height-1 {
		return true
	}
	// Check if touch itself.
	// Don't check tail
	for i := 1; i < len(ss.SnakeBody); i++ {
		body := ss.SnakeBody[i].Cord
		if new_head.X == body.X && new_head.Y == body.Y {
			return true
		}
	}

	return false
}

func (ss *SnakeState) maybe_consume_apple_and_grow_snake(new_head SnakePart) {
	// grow snake to new head
	ss.SnakeBody = append(ss.SnakeBody, new_head)
	if new_head.Cord == ss.Apple {
		// consume apple
		ss.Apple = Point{-1, -1}
		ss.Score += 1
	} else {
		// cut off tail, this will make the snake move one cell
		ss.SnakeBody = ss.SnakeBody[1:]
	}

	if len(ss.SnakeBody) > 1 {
		// More than 1 body, the last one is tail
		tail_type, err := get_tail_type(ss.SnakeBody[0].Cord, ss.SnakeBody[1].Cord)
		if err == nil {
			ss.SnakeBody[0].PartType = tail_type
		} else {
			log.Fatal(err)
		}
	}

	if len(ss.SnakeBody) > 2 {
		// More than 2, there might be turns, only need to
		// update the body type of the old head, that's where
		// the turn happens
		t, err := get_part_type(
			ss.SnakeBody[len(ss.SnakeBody)-1].Cord,
			ss.SnakeBody[len(ss.SnakeBody)-2].Cord,
			ss.SnakeBody[len(ss.SnakeBody)-3].Cord)
		if err == nil {
			ss.SnakeBody[len(ss.SnakeBody)-2].PartType = t
		} else {
			log.Fatal(err)
		}
	}
}

var (
	_BODY_TYPE_DELTAS = map[[4]int8]int{
		// X are all the same, the body is straight up or down
		{0, 0, 1, 1}: BODY_PART_I,
		// Y are all the same, the body is straight left or right
		{1, 1, 0, 0}: BODY_PART_H,
		// L shape
		{0, 1, 1, 0}: BODY_PART_BODY_L,
		// L1, L rotated 90 clockwise
		{0, 1, -1, 0}: BODY_PART_BODY_L1,
		// L2, L rotated 180 clock wise
		{0, -1, -1, 0}: BODY_PART_BODY_L2,
		// L3, L rotated 270 clock wise
		{0, -1, 1, 0}: BODY_PART_BODY_L3,
	}

	_TAIL_TYPE_DELTAS = map[[2]int8]int{
		{0, 1}:  BODY_PART_TAIL_DOWN,
		{0, -1}: BODY_PART_TAIL_UP,
		{-1, 0}: BODY_PART_TAIL_LEFT,
		{1, 0}:  BODY_PART_TAIL_RIGHT,
	}

	_HEAD_TYPE_FROM_DIR = map[int]int{
		UP:    BODY_PART_HEAD_UP,
		DOWN:  BODY_PART_HEAD_DOWN,
		LEFT:  BODY_PART_HEAD_LEFT,
		RIGHT: BODY_PART_HEAD_RIGHT,
	}
)

// Returns the body type of p2
func get_part_type(p1, p2, p3 Point) (int, error) {
	delta1 := [4]int8{
		int8(p2.X - p1.X),
		int8(p3.X - p2.X),
		int8(p2.Y - p1.Y),
		int8(p3.Y - p2.Y)}
	delta2 := [4]int8{
		int8(p2.X - p3.X),
		int8(p1.X - p2.X),
		int8(p2.Y - p3.Y),
		int8(p1.Y - p2.Y)}
	type1, type1_ok := _BODY_TYPE_DELTAS[delta1]
	type2, type2_ok := _BODY_TYPE_DELTAS[delta2]
	if type1_ok {
		return type1, nil
	} else if type2_ok {
		return type2, nil
	} else {
		return 0, fmt.Errorf("unknow body type %v, %v, %v", p1, p2, p3)
	}
}

func get_tail_type(p_tail, p_pre Point) (int, error) {
	delta := [2]int8{int8(p_pre.X - p_tail.X), int8(p_pre.Y - p_tail.Y)}
	type1, ok := _TAIL_TYPE_DELTAS[delta]
	if ok {
		return type1, nil
	} else {
		return 0, fmt.Errorf("unknown tail type %v, %v", p_tail, p_pre)
	}
}
