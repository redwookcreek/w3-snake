package snake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnakeDirection(t *testing.T) {
	snake_state := CreateSnake(10, 10)

	assert.Equal(t, DOWN, snake_state.Direction)
	assert.False(t, snake_state.GameOver)

	// Cannot reverse the snake
	snake_state.UpdateDirection(UP)
	assert.Equal(t, DOWN, snake_state.Direction)

	// same direction
	snake_state.UpdateDirection(DOWN)
	assert.Equal(t, DOWN, snake_state.Direction)

	// Update direction
	snake_state.UpdateDirection(LEFT)
	assert.Equal(t, LEFT, snake_state.Direction)
	snake_state.UpdateDirection(RIGHT)
	assert.Equal(t, LEFT, snake_state.Direction)

	snake_state.UpdateDirection(DOWN)
	assert.Equal(t, DOWN, snake_state.Direction)

}

func TestCreateApple(t *testing.T) {
	snake_state := CreateSnake(10, 10)
	assert.False(t, snake_state.HasApple())

	// Initial snake has no apple
	assert.Less(t, snake_state.Apple.X, 0)
	assert.Less(t, snake_state.Apple.Y, 0)

	snake_state.maybe_create_apple()
	assert.True(t, snake_state.HasApple())

	// randomly created apple should be in board
	assert.Greater(t, snake_state.Apple.X, 0)
	assert.Less(t, snake_state.Apple.X, snake_state.Width)
	assert.Greater(t, snake_state.Apple.Y, 0)
	assert.Less(t, snake_state.Apple.Y, snake_state.Height)
}

func make_head(x, y, part int) SnakePart {
	return SnakePart{
		Point{x, y},
		part,
	}
}

func make_body(x, y int) SnakePart {
	return SnakePart{
		Point{x, y},
		BODY_PART_I,
	}
}

func make_tail(x, y, part int) SnakePart {
	return SnakePart{
		Point{x, y},
		part,
	}
}
func TestAdvanceSnakeHead(t *testing.T) {
	snake_state := CreateSnake(10, 10)

	// Initial snake with just one cell, at the center of the board
	// moving down ward
	assert.Equal(t, []SnakePart{make_head(5, 5, BODY_PART_HEAD_DOWN)}, snake_state.SnakeBody)
	assert.Equal(t, DOWN, snake_state.Direction)

	assert.Equal(t, make_head(5, 6, BODY_PART_HEAD_DOWN), snake_state.advance_snake_head())

	snake_state.UpdateDirection(LEFT)
	assert.Equal(t, LEFT, snake_state.Direction)
	assert.Equal(t, make_head(4, 5, BODY_PART_HEAD_LEFT), snake_state.advance_snake_head())

	snake_state.UpdateDirection(UP)
	assert.Equal(t, UP, snake_state.Direction)
	assert.Equal(t, make_head(5, 4, BODY_PART_HEAD_UP), snake_state.advance_snake_head())

	snake_state.UpdateDirection(RIGHT)
	assert.Equal(t, make_head(6, 5, BODY_PART_HEAD_RIGHT), snake_state.advance_snake_head())
}

func TestSnakeTouched(t *testing.T) {
	snake_state := CreateSnake(10, 10)

	// not touched
	assert.False(t, snake_state.snake_touched(Point{4, 5}))
	assert.False(t, snake_state.snake_touched(Point{3, 5}))

	// touch left
	assert.True(t, snake_state.snake_touched(Point{0, 5}))
	assert.True(t, snake_state.snake_touched(Point{-1, 5}))

	// touch right
	assert.True(t, snake_state.snake_touched(Point{10, 5}))
	assert.True(t, snake_state.snake_touched(Point{11, 5}))

	// touch top
	assert.True(t, snake_state.snake_touched(Point{5, 0}))
	assert.True(t, snake_state.snake_touched(Point{5, -1}))

	// touch bottom
	assert.True(t, snake_state.snake_touched(Point{5, 10}))
	assert.True(t, snake_state.snake_touched(Point{5, 11}))

	// Touch self
	snake_state.SnakeBody = []SnakePart{
		make_head(4, 5, BODY_PART_HEAD_DOWN), make_body(5, 5), make_tail(6, 5, BODY_PART_TAIL_DOWN)}
	assert.False(t, snake_state.snake_touched(Point{7, 7}))
	assert.True(t, snake_state.snake_touched(Point{5, 5}))
	// touching tail does not count
	assert.False(t, snake_state.snake_touched(Point{4, 5}))
}

func TestConsumeAppleAndGrowSnake(t *testing.T) {
	// snake head and apple are both at {5, 5}
	ss := CreateSnake(10, 10)
	ss.Apple = Point{5, 6}
	// snake should grow, and apple should disappear
	ss.maybe_consume_apple_and_grow_snake(make_head(5, 6, BODY_PART_HEAD_DOWN))
	assert.Equal(t, Point{-1, -1}, ss.Apple)
	assert.Equal(t, []SnakePart{
		make_tail(5, 5, BODY_PART_TAIL_DOWN),
		make_head(5, 6, BODY_PART_HEAD_DOWN)}, ss.SnakeBody)
	assert.Equal(t, 1, ss.Score)

	// Snake head does not touch apple
	ss = CreateSnake(10, 10)
	ss.Apple = Point{2, 3}
	ss.maybe_consume_apple_and_grow_snake(make_head(5, 6, BODY_PART_HEAD_DOWN))
	// Apple remains the same
	// snake moved to new head
	assert.Equal(t, Point{2, 3}, ss.Apple)
	assert.Equal(t, []SnakePart{make_head(5, 6, BODY_PART_HEAD_DOWN)}, ss.SnakeBody)
}

func TestTick(t *testing.T) {
	ss := CreateSnake(10, 10)
	ss.Apple = Point{5, 6}
	ss.Tick()
	assert.Equal(t, []SnakePart{
		make_tail(5, 5, BODY_PART_TAIL_DOWN),
		make_head(5, 6, BODY_PART_HEAD_DOWN)}, ss.SnakeBody)
	assert.NotEqual(t, Point{-1, -1}, ss.Apple) // should generate a new apple

	ss.UpdateDirection(LEFT)
	ss.Tick()
	assert.Equal(t, []SnakePart{
		make_tail(5, 6, BODY_PART_TAIL_LEFT), make_head(4, 6, BODY_PART_HEAD_LEFT)}, ss.SnakeBody)

	for i := 0; i < 7; i++ {
		ss.Tick()
	}
	// Should touched wall
	assert.True(t, ss.GameOver)
	// Snake should remain at the location before touching the wall
	assert.Equal(t, make_head(0, 6, BODY_PART_HEAD_LEFT), ss.advance_snake_head())
}

func assert_body_type(t *testing.T, p1, p2, p3 Point, expected_body_type int) {
	type1, _ := get_part_type(p1, p2, p3)
	assert.Equal(t, expected_body_type, type1)
	// Reverse the point sequence should get the same body type
	type2, _ := get_part_type(p3, p2, p1)
	assert.Equal(t, expected_body_type, type2)
}
func TestGetBodyPart(t *testing.T) {
	assert_body_type(
		t,
		Point{4, 5}, Point{4, 6}, Point{4, 7},
		BODY_PART_I)
	assert_body_type(
		t,
		Point{4, 5}, Point{5, 5}, Point{6, 5},
		BODY_PART_H)
	assert_body_type(
		t,
		Point{4, 5}, Point{4, 6}, Point{5, 6},
		BODY_PART_BODY_L)
	assert_body_type(
		t,
		Point{5, 6}, Point{5, 5}, Point{6, 5},
		BODY_PART_BODY_L1)
	assert_body_type(
		t,
		Point{4, 5}, Point{5, 5}, Point{5, 6},
		BODY_PART_BODY_L2)
	assert_body_type(
		t,
		Point{4, 5}, Point{5, 5}, Point{5, 4},
		BODY_PART_BODY_L3)
}
