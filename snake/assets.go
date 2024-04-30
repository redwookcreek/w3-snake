package snake

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	//go:embed snake-graphics.png
	asset_sprite_file []byte

	//go:embed *.png
	sset_folder embed.FS

	// Font face for display game info
	M_PLUS_FACE_SCOURCE *text.GoTextFaceSource

	// Image assets
	ASSETS_SPRITE *ebiten.Image // This is the full sprite

	// Subimage of the sprite
	SNAKE_HEAD_UP_IMG    *ebiten.Image
	SNAKE_HEAD_DOWN_IMG  *ebiten.Image
	SNAKE_HEAD_LEFT_IMG  *ebiten.Image
	SNAKE_HEAD_RIGHT_IMG *ebiten.Image

	SNAKE_TAIL_UP_IMG    *ebiten.Image
	SNAKE_TAIL_DOWN_IMG  *ebiten.Image
	SNAKE_TAIL_LEFT_IMG  *ebiten.Image
	SNAKE_TAIL_RIGHT_IMG *ebiten.Image

	SNAKE_BODY_HORIZONTAL_IMG *ebiten.Image
	SNAKE_BODY_VERTICAL_IMG   *ebiten.Image

	SNAKE_BODY_L_IMG  *ebiten.Image
	SNAKE_BODY_L1_IMG *ebiten.Image
	SNAKE_BODY_L2_IMG *ebiten.Image
	SNAKE_BODY_L3_IMG *ebiten.Image

	APPLE_IMG *ebiten.Image

	BODY_PART_TO_IMG_MAP map[int]*ebiten.Image

	BOARDER_COLOR = color.White

	sprite_cell_size = 320 / 5
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	M_PLUS_FACE_SCOURCE = s

	ASSETS_SPRITE, _, err = ebitenutil.NewImageFromFile("snake/snake-graphics.png")
	if err != nil {
		log.Fatal(err)
	}

	SNAKE_HEAD_UP_IMG = sprite_sub_image(0, 3)
	SNAKE_HEAD_RIGHT_IMG = sprite_sub_image(0, 4)
	SNAKE_HEAD_DOWN_IMG = sprite_sub_image(1, 4)
	SNAKE_HEAD_LEFT_IMG = sprite_sub_image(1, 3)

	SNAKE_TAIL_UP_IMG = sprite_sub_image(2, 3)
	SNAKE_TAIL_DOWN_IMG = sprite_sub_image(3, 4)
	SNAKE_TAIL_LEFT_IMG = sprite_sub_image(3, 3)
	SNAKE_TAIL_RIGHT_IMG = sprite_sub_image(2, 4)

	SNAKE_BODY_HORIZONTAL_IMG = sprite_sub_image(0, 1)
	SNAKE_BODY_VERTICAL_IMG = sprite_sub_image(1, 2)
	SNAKE_BODY_L_IMG = sprite_sub_image(1, 0)
	SNAKE_BODY_L1_IMG = sprite_sub_image(0, 0)
	SNAKE_BODY_L2_IMG = sprite_sub_image(0, 2)
	SNAKE_BODY_L3_IMG = sprite_sub_image(2, 2)

	APPLE_IMG = sprite_sub_image(3, 0)

	BODY_PART_TO_IMG_MAP = map[int]*ebiten.Image{
		BODY_PART_HEAD_DOWN:  SNAKE_HEAD_DOWN_IMG,
		BODY_PART_HEAD_UP:    SNAKE_HEAD_UP_IMG,
		BODY_PART_HEAD_LEFT:  SNAKE_HEAD_LEFT_IMG,
		BODY_PART_HEAD_RIGHT: SNAKE_HEAD_RIGHT_IMG,

		BODY_PART_TAIL_DOWN:  SNAKE_TAIL_DOWN_IMG,
		BODY_PART_TAIL_UP:    SNAKE_TAIL_UP_IMG,
		BODY_PART_TAIL_LEFT:  SNAKE_TAIL_LEFT_IMG,
		BODY_PART_TAIL_RIGHT: SNAKE_TAIL_RIGHT_IMG,

		BODY_PART_I:       SNAKE_BODY_VERTICAL_IMG,
		BODY_PART_H:       SNAKE_BODY_HORIZONTAL_IMG,
		BODY_PART_BODY_L:  SNAKE_BODY_L_IMG,
		BODY_PART_BODY_L1: SNAKE_BODY_L1_IMG,
		BODY_PART_BODY_L2: SNAKE_BODY_L2_IMG,
		BODY_PART_BODY_L3: SNAKE_BODY_L3_IMG,
	}
}

func sprite_sub_image(row, col int) *ebiten.Image {
	return ASSETS_SPRITE.SubImage(sprite_rect(row, col)).(*ebiten.Image)
}

// Give a row, col of the sprite, return
// the rectangle that contains the sub image of the sprite
func sprite_rect(row, col int) image.Rectangle {
	// The sprite is a 320 by 256 pixel image
	// with 4 rows and 5 columes
	min_point := image.Point{
		col * sprite_cell_size,
		row * sprite_cell_size}
	max_point := image.Point{
		min_point.X + sprite_cell_size,
		min_point.Y + sprite_cell_size}
	return image.Rectangle{min_point, max_point}
}

// Return the scale factor for drawing the cell on screen
func CellScale(cell_size float64) float64 {
	return cell_size / float64(sprite_cell_size)
}
