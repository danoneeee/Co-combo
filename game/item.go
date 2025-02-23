package game

import resource "github.com/quasilyte/ebitengine-resource"

type Item struct {
	X, Y             int     // Позиция изображения
	Dragging         bool    // Флаг, указывающий, что изображение перемещается
	OffsetX, OffsetY float64 // Смещение курсора относительно изображения
	TypeObject       resource.ImageID
}
